package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"	

	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/logger"
	"github.com/dfklegend/cell/utils/runservice"	

	rpccommon "github.com/dfklegend/cell/rpc/common"
	"github.com/dfklegend/cell/rpc/consts"
	"github.com/dfklegend/cell/rpc/interfaces"
	"github.com/dfklegend/cell/rpc/config"
	"github.com/dfklegend/cell/rpc/init"
	"github.com/dfklegend/cell/rpc/protos"	
)

func init() {
	initp.Visit()
}

var debugFlags = rpccommon.GetDebugFlags()
var MailBoxIdService *common.SerialIdService = common.NewSerialIdService()

type MailBoxStat struct {

	// 调用的请求
	ReqCalled uint64
	//
	ReqNeedAckCalled uint64
	ReqAckGot        uint64
}

// 邮箱，代表一个可访问地址 
// pending实际逻辑都被Post回到了MailBox主routine执行
// 主routine
// 		callrequest
// 		checkPending
// CheckPendingTimer
// 		post checkPending调用	
// 		
// pull ack timer
// 		PullAck，结果post到主routine
// 		post 检查ack过期

// 减少异步情况
type MailBox struct {
	MailBoxId uint32
	Address   string
	// 所属station
	station *MailStation
	
	connectCalled bool
	impl interfaces.IRPCClientImpl

	// 缓存请求队列
	// 当连接还未成功之时
	pendings map[rpccommon.ReqIdType]*rpccommon.Request
	// 等待返回的请求
	waitAcks map[rpccommon.ReqIdType]*rpccommon.Request

	// 
	mutex sync.Mutex
	// 连续空拉
	conEmptyPullCount int32

	// 尝试建一个routine来处理
	runService         *runservice.RunService	
	
	wgStop 				sync.WaitGroup
	stopped bool

	// 不考虑同步
	Stat MailBoxStat
}

func NewMailBox() *MailBox {
	id := MailBoxIdService.AllocId()
	return &MailBox{
		MailBoxId: id,
		pendings:  make(map[rpccommon.ReqIdType]*rpccommon.Request),
		waitAcks:  make(map[rpccommon.ReqIdType]*rpccommon.Request),		
		stopped:   false,
		connectCalled: false,
	}
}

func (self *MailBox) SetStation(station *MailStation) {
	self.station = station
}

func (self *MailBox) GetStationIId() string {
	if self.station != nil {
		return self.station.GetStationId()
	}
	return "notset"
}

func (self *MailBox) GetRunService() *runservice.RunService {
	return self.runService
}

// 异步连接
// lazy连接
func (self *MailBox) Start(address string) {

	impl := config.CreateClientImpl()
	impl.SetMailBox(self)
	self.impl = impl

	self.runService = runservice.NewRunService(fmt.Sprintf("mailbox%v", self.MailBoxId))
	self.runService.Start()
	self.Address = address	

	// 加入一个检查pending的timer	
	self.startCheckPendings()

	// 定时拉取rpc返回
	self.startPullAcks()
	self.wgStop.Add(2)
}

// TODO: 如何保证安全关闭
func (self *MailBox) Stop() {
	if self.impl == nil {
		return
	}	
	
	go func(){
		// 等待子routine结束
		self.wgStop.Wait()
		//self.mutex.Lock()
		//defer self.mutex.Unlock()
		self.impl.Close()
		self.runService.Stop()
	}()
	//self.runService.Stop()
	//self.conn.Close()	
	self.stopped = true
}

func (self *MailBox) tryLazyConnect() {
	if self.connectCalled {
		return
	}
	self.connectCalled = true;
	self.doConnect(self.Address)
}

// 异步连接
func (self *MailBox) doConnect(address string) {	
	self.impl.Connect(address)
}

func (self *MailBox) startCheckPendings() {
	t := time.NewTicker(1 * time.Second)	

	go func() {
		defer t.Stop()

		for !self.stopped {
			if self.isConnected() && len(self.pendings) == 0 {
				// 简单sleep释放点cpu
				time.Sleep(1 * time.Second)
			}			

			select {
			// case <-self.stopChan:
			// 	return
			case <-t.C:
				self.checkPendings()
			}
		}

		self.wgStop.Done()
	}()
}

func (self *MailBox) needPending() bool {
	return !self.isConnected() || self.hasPendings()
}

func (self *MailBox) hasPendings() bool {
	//self.mutex.Lock()
	//defer self.mutex.Unlock()
	return len(self.pendings) > 0
}

// 
func (self *MailBox) checkPendings() {
	//logger.Log.Debugf("checkPendings")
	self.runService.GetScheduler().Post(func() (interface{}, error) {
		self.doCheckPendings()
		return nil, nil
	})
}

// 执行pending
// 找到最小的key
// 
// TODO: 检查pending超时
func (self *MailBox) doCheckPendings() {
	if len(self.pendings) == 0 || !self.isConnected() {
		return
	}

	checkTimes := 1000
	half := len(self.pendings) / 2
	if checkTimes < half {
		checkTimes = half
	}

	for i := 0; i < checkTimes; i++ {
		head := self.popHead()
		if head == nil || self.stopped {
			return
		}
		logger.Log.Infof("run pending:%v\n", head.ReqId)
		self.runRequest(head)
	}
}

func (self *MailBox) popHead() *rpccommon.Request {
	// 找到最小的key
	//self.mutex.Lock()
	//defer self.mutex.Unlock()

	var head *rpccommon.Request
	for id, one := range self.pendings {
		if head == nil {
			head = one
		} else {
			if id < head.ReqId {
				head = one
			}
		}
	}

	if head != nil {
		delete(self.pendings, head.ReqId)
	}
	return head
}

// 后续可以调整频率，比如，如果长时间没有拉到数据
// 增加拉取频率，发送消息后，立刻重置
func (self *MailBox) startPullAcks() {
	// 实测30ms效果较好
	t := time.NewTicker(30 * time.Millisecond)

	go func() {
		defer t.Stop()
		for !self.stopped {
			// 拉取和等待检查
			if self.conEmptyPullCount > 100 {
				// 释放点CPU
				//logger.Log.Debugf("pull sleep")
				time.Sleep(1 * time.Second)
			}

			select {
			// case <-self.stopChan:
			// 	return
			case <-t.C:
				self.updateAcks()
			}
		}

		self.wgStop.Done()
	}()
}

func (self *MailBox) updateAcks() {
	self.pullAcks()

	//logger.Log.Debugf("updateAcks")
	self.runService.GetScheduler().Post(func() (interface{}, error) {
		self.processWaitAckTimeout()
		return nil, nil
	})		
}

func (self *MailBox) pullAcks() {	
	self.impl.ReqPullAcks()
}

func (self *MailBox) OnGotAcksCount(acksCount int32) {
	if acksCount == 0 {
		self.conEmptyPullCount	++
	} else {
		self.conEmptyPullCount = 0
	}
	
}

// 多久没收到反馈的rpc，自动ack timeout
func (self *MailBox) processWaitAckTimeout() {
	if len(self.waitAcks) == 0 {
		return
	}

	// 循环查找
	now := time.Now().Unix()
	expired := make([]rpccommon.ReqIdType, 0)

	//self.mutex.Lock()
	for id, one := range self.waitAcks {
		if now >= one.TimeOut {
			// need remove
			expired = append(expired, id)
		}
	}
	//self.mutex.Unlock()

	if len(expired) == 0 {
		return
	}

	errStr := consts.ErrRPCAutoTimeout.Error()

	for _, id := range expired {
		self.ApplyAck(uint32(id), errStr, "")
		self.RemoveAck(uint32(id))

		logger.Log.Errorf("req %v timeout", id)
	}
}

func (self *MailBox) isConnected() bool {	
	return self.impl.IsConnected()
}

// context过期判定
// cb(data, error)
func (self *MailBox) Call(method string,
	inArg interface{}, options ...rpccommon.ReqOption) {

	if self.stopped {
		logger.Log.Errorf("Call to a mailbox stopped!")
		return
	}

	// 不锁 极端情况可能会报post chan关闭了
	// 暂时不锁
	//self.mutex.Lock()
	//defer self.mutex.Unlock()
	self.runService.GetScheduler().Post(func() (interface{}, error) {
		self.doCall(method, inArg, options...)
		return nil, nil
	})	
}

// 未来可以考虑能并发发起请求
func (self *MailBox) doCall(method string,
	inArg interface{}, options ...rpccommon.ReqOption) {

	self.tryLazyConnect()
	// 构建request
	req := rpccommon.NewRequest(method, inArg, options...)

	// 如果发现连接未建立，push到pending
	// 如果要保证有序，存在pending时，也先push
	if self.needPending() {
		self.pushPending(req)
		return
	}

	self.runRequest(req)
}

func (self *MailBox) pushPending(req *rpccommon.Request) {
	//logger.Log.Infof("push Pending:%v\n", req.ReqId)
	if len(self.pendings) > 100000 {
		logger.Log.Errorf("TOO MANY PENDING REQ! skip it\n")
		return
	}
	self.pushReqMap(self.pendings, req)
}

func (self *MailBox) pushWaitAck(req *rpccommon.Request) {
	self.pushReqMap(self.waitAcks, req)
}

func (self *MailBox) pushReqMap(m map[rpccommon.ReqIdType]*rpccommon.Request, req *rpccommon.Request) {
	//self.mutex.Lock()
	//defer self.mutex.Unlock()

	if m[req.ReqId] != nil {
		logger.Log.Infof("duplicate req:%v", req.ReqId)
		return
	}

	m[req.ReqId] = req
}

// run in runservice
func (self *MailBox) runRequest(req *rpccommon.Request) {
	self.callRequest(req)
	self.Stat.ReqCalled++
	if req.NeedAck() {
		// push to acks
		self.pushWaitAck(req)
		self.Stat.ReqNeedAckCalled++
	}
}

func (self *MailBox) callRequest(req *rpccommon.Request) {	
	if debugFlags.RPCNotSend {
		return
	}
	self.doRequest(req)
}

func (self *MailBox) doRequest(req *rpccommon.Request) {
	self.impl.CallRPC(req)	
}

func (self *MailBox) getAck(reqId uint32) *rpccommon.Request {
	//self.mutex.Lock()
	//defer self.mutex.Unlock()
	return self.waitAcks[rpccommon.ReqIdType(reqId)]
}

func (self *MailBox) RemoveAck(reqId uint32) {
	//self.mutex.Lock()
	//defer self.mutex.Unlock()

	req := self.waitAcks[rpccommon.ReqIdType(reqId)]
	if req != nil {
		delete(self.waitAcks, rpccommon.ReqIdType(reqId))
	}
}

// . 查找等待的req是否还在(可能已经timeout出去了)
// . 还存在 根据返回类型，生成回调数据
// . 调用回调函数
func (self *MailBox) ApplyAck(reqId uint32, errStr string, message string) {
	self.Stat.ReqAckGot++

	req := self.getAck(reqId)
	if req == nil {
		logger.Log.Errorf("ack can not find req:%v\n", reqId)
		return
	}

	args := MakeCBArgs(errStr, message, req.CBArgType)
	//logger.Log.Infof("args:%v", args)
	if req.Scheduler != nil {
		req.Scheduler.Post(func() (interface{}, error) {
			req.CBFunc.Call(args)
			return nil, nil
		})
	} else {
		req.CBFunc.Call(args)
	}
}

func MakeCBArgs(errStr string, message string, cbArgType reflect.Type) []reflect.Value {
	var argValue, errValue reflect.Value

	// make data
	data := reflect.New(cbArgType.Elem()).Interface()
	if errStr == "" {
		errJson := json.Unmarshal([]byte(message), data)
		//logger.Log.Infof("data:%v", data)
		if errJson != nil {
			logger.Log.Infof("error parse JSON:%v", message)
		}
		errValue = consts.NilError
	} else {
		errValue = reflect.ValueOf(errors.New(errStr))
	}
	argValue = reflect.ValueOf(data)

	args := []reflect.Value{argValue, errValue}
	return args
}

func callRPC(ctx context.Context, c protos.RPCServerClient,
	method string, message string) (*protos.RPCReply, error) {
	logger.Log.Infof("pre call callRPC routine:%v", common.GetRoutineID())
	r, err := c.RPC(ctx, &protos.RPCRequest{Method: method, Message: message})
	if err != nil {
		logger.Log.Fatalf("could not call rpc: %v", err)
	}
	logger.Log.Infof("post routine:%v Greeting: %v", common.GetRoutineID(), r)
	return r, err
}

func (self *MailBox) DumpStat() {
	stat := self.Stat
	logger.Log.Infof("---- stat ----\n")
	logger.Log.Infof("ReqCalled:%v\n", stat.ReqCalled)
	logger.Log.Infof("ReqNeedAckCalled:%v\n", stat.ReqNeedAckCalled)
	logger.Log.Infof("ReqAckGot:%v\n", stat.ReqAckGot)
}

func (self *MailBox) IsBenchOver() bool {
	stat := self.Stat

	return stat.ReqNeedAckCalled > 0 && stat.ReqNeedAckCalled == stat.ReqAckGot
}
