package server

import (
	"log"

	api "github.com/dfklegend/cell/rpc/apientry"
	rpccommon "github.com/dfklegend/cell/rpc/common"
	"github.com/dfklegend/cell/rpc/protos"
	"github.com/dfklegend/cell/rpc/interfaces"
	"github.com/dfklegend/cell/rpc/config"
	"github.com/dfklegend/cell/rpc/init"
	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/runservice"
)

func init() {
	initp.Visit()
}

type RPCServerNode struct {
	Name string
	// 接口集
	collection *api.APICollection
	// 对应的rpc 适配器
	impl     interfaces.IRPCServerImpl
	waitAckCtrl *rpccommon.WaitAckCtrl

	runService *runservice.RunService

	reqIdService *common.SerialIdService
}

func NewNode(name string) *RPCServerNode {
	return &RPCServerNode{
		Name:         name,
		collection:   api.NewCollection(),
		impl:      nil,
		waitAckCtrl:  rpccommon.NewWaitAckCtrl(),
		runService:   runservice.NewRunService(name),
		reqIdService: common.NewSerialIdService(),
	}
}

func (self *RPCServerNode) GetCollection() *api.APICollection {
	return self.collection
}

// TODO: 配置对象
func (self *RPCServerNode) Start(listen string) {
	impl := config.CreateServerImpl()
	impl.Init(self)
	impl.Start(listen)
	self.impl = impl

	self.collection.Build()
}

func (self *RPCServerNode) Stop() {
	if self.impl != nil {
		self.impl.Stop()
	}
}

func (self *RPCServerNode) makeWaitAckFromReq(in *protos.RPCRequest) *rpccommon.WaitAckReq {
	wa := rpccommon.NewWaitAckReq()
	wa.ServerReqId = rpccommon.ReqIdType(self.reqIdService.AllocId())
	wa.ClientId = in.ClientId
	wa.ClientReqId = rpccommon.ReqIdType(in.ReqId)
	return wa
}

func (self *RPCServerNode) Call(in *protos.RPCRequest) {
	// 注册等待Ack
	// 调度 执行
	// push ack
	needAck := in.NeedAck
	var serverReqId rpccommon.ReqIdType
	serverReqId = 0

	if needAck {
		wa := self.makeWaitAckFromReq(in)
		self.waitAckCtrl.RegisterWaitAck(wa)
		serverReqId = wa.ServerReqId
	}

	//log.Printf("node.call:%v\n", in)
	callerr := self.collection.Call(in.Method, []byte(in.Message), func(e error, result interface{}) {
		//log.Printf("node call got:%v %vn", e, result)
		if !needAck {
			return
		}

		errStr := ""
		message := ""

		if e != nil {
			errStr = e.Error()
		} else {
			message = string(result.([]byte))
		}

		// 等着客户端来拉取
		self.waitAckCtrl.PushAck(serverReqId, errStr, message)
	}, nil)

	// 调用错误
	if callerr != nil {
		log.Printf("callerr:%v\n", callerr)
		if needAck {
			self.waitAckCtrl.PushAck(serverReqId, callerr.Error(), "")
			//self.waitAckCtrl.RemoveWaitAck(serverReqId)
		}
	}
}

func (self *RPCServerNode) PopReadyAcks(clientId string, maxNum int) []*rpccommon.WaitAckReq {
	return self.waitAckCtrl.PopReadyAcks(clientId, maxNum)
}

func (self *RPCServerNode) IsTooBusy() bool {
	// TODO: 发现等待返回的结果过多,返回true
	return false
}
