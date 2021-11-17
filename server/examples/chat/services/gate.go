package services

import (	
	"strings"
	"strconv"
	"reflect"
	"log"
	"fmt"
	"math/rand"

	api "github.com/dfklegend/cell/rpc/apientry"		
	rpccommon "github.com/dfklegend/cell/rpc/common"
	"github.com/dfklegend/cell/utils/logger"	
	"github.com/dfklegend/cell/utils/waterfall"	
	ucommon "github.com/dfklegend/cell/utils/common"	

	"github.com/dfklegend/cell/server/common"
	"github.com/dfklegend/cell/server/client/session"
	"github.com/dfklegend/cell/server/interfaces"
	"github.com/dfklegend/cell/server/examples/chat/protos"
)

var (
	tempIdService *ucommon.SerialIdService = ucommon.NewSerialIdService()
)

func init() {
	session.AddOnSessionClose(onSessionClose)
}

func onSessionClose(s session.IServerSession) {
	app := interfaces.GetApp()

	log.Printf("%+v\n", s)
	//fs := s.(*session.FrontSession)
	//log.Printf("%+v %v %v\n", fs.Data, fs.Get("chatid", 0), s.Get("chatid", 0))
	s.Lock()
	uid := s.GetID()
	frontId := session.GetServerId(s)
	netId := session.GetNetId(s)	
	chatId := s.Get("chatid", 0)
	serverId := fmt.Sprintf("chat-%v", chatId)
	log.Printf("%v\n", serverId)
	s.Unlock()

	// 没有绑定服务器
	if chatId == 0 {
		return
	}


	ms := app.GetMailStation()
	go func(){
		waterfall.Waterfall_Go([]waterfall.Task{func(args ...interface{}) {
			callback, _ := args[0].(waterfall.Callback)
			
			ms.Call(serverId, "chatremote.leave", &protos.RoomLeave{
					UId: uid,
					ServerId: frontId,
					NetId: netId,
				}, 
				rpccommon.ReqWithCB(reflect.ValueOf(func(result *protos.NormalAck, e error) {
					log.Printf("rpc got result:%v %v\n", result, e)
					if e != nil {
						log.Printf("callback(true)\n")
						callback(true, e)	
					} else {
						callback(false, nil)	
					}
			})))
		}}, func(args ...interface{}) {
			isErr, _ := args[0].(bool)
			log.Printf("args:%+v\n", args)
			log.Printf("isErr:%v\n", isErr)			
		})
	}()
}

type Gate struct {
	api.APIEntry
}

// 分配一个具体连接的gate
func (self *Gate) QueryGate(d api.IHandlerSession, msg *protos.EmptyArgReq, cbFunc api.HandlerCBFunc) error {
	logger.Log.Debugf("enter gate.QueryGate")

	//fs := d.(*session.FrontSession)

	app := interfaces.GetApp()
	
	servers, _ := app.GetServersByType("gate")	

	logger.Log.Debugf("servers:%+v", servers)
	if servers == nil || len(servers) == 0 {
		return nil
	}

	var server *common.Server

	var wantIndex = rand.Intn(len(servers))
	index := 0
	for _, v := range(servers) {
		if index == wantIndex {
			server = v
			break
		}
		index ++
	}

	address := server.WSClientAddress

	ip := ""
	port1 :=  ""
	port2 :=  ""
	subs := strings.Split(address, ":")
	if len(subs) == 2 {
		ip = subs[0]
		port1 = subs[1]
	}

	address = server.ClientAddress
	subs = strings.Split(address, ":")
	if len(subs) == 2 {		
		port2 = subs[1]
	}
	
	api.CheckInvokeCBFunc(cbFunc, nil,
		&protos.QueryGateAck{
			Code: 0,
			IP: ip,
			Port: fmt.Sprintf("%v,%v", port1, port2),
			})
	return nil
}

func (self *Gate) doLoginRet(cbFunc api.HandlerCBFunc, code int, errStr string) {
	api.CheckInvokeCBFunc(cbFunc, nil,
		&protos.NormalAck{
		Code: code,
		Result: errStr,
	})	
}

// 登录
// 分配并绑定到某个chat服务器
func (self *Gate) Login(d api.IHandlerSession, msg *protos.LoginReq, cbFunc api.HandlerCBFunc) error {
	fs := d.(*session.FrontSession)

	app := interfaces.GetApp()	
	servers, _ := app.GetServersByType("chat")	

	logger.Log.Debugf("login:%+v", msg)
	logger.Log.Debugf("servers:%+v", servers)
	if servers == nil || len(servers) == 0 {
		self.doLoginRet(cbFunc, 1, "no chat server")
		return nil
	}

	// 随机选一个chat
	var serverId = ""
	var wantIndex = rand.Intn(len(servers))
	//logger.Log.Debugf("wantIndex:%v\n", wantIndex)
	index := 0
	for k, _ := range(servers) {
		if index == wantIndex {
			serverId = k
			break
		}
		index ++		
	}
	logger.Log.Debugf("goto:%v\n", serverId)

	subs := strings.Split(serverId, "-")
	if len(subs) < 2 {
		self.doLoginRet(cbFunc, 1, "unknown err")
		return nil
	}

	index, _ = strconv.Atoi(subs[1])

	uid := fmt.Sprintf("uid-%v", tempIdService.AllocId())

	fs.Lock()
	fs.Bind(uid)
	fs.Set("chatid", index)
	fs.PushSession()
	fs.Unlock()


	// test
	fs.Lock()
	frontId := session.GetServerId(fs)
	netId := session.GetNetId(fs)
	fs.Unlock()
	app.GetSysService().PushMessageById(frontId, netId, "onTest", "{}", nil)
	//
	

	// 注册onClose

	// 通知对方chat服务器，有玩家进入
	ms := app.GetMailStation()
	go func(){
		waterfall.Waterfall_Go([]waterfall.Task{func(args ...interface{}) {
			callback, _ := args[0].(waterfall.Callback)
			
			ms.Call(serverId, "chatremote.entry", &protos.RoomEntry{
					UId: uid, 
					Name: msg.Name,
					ServerId: frontId,
					NetId: netId,
				}, 
				rpccommon.ReqWithCB(reflect.ValueOf(func(result *protos.NormalAck, e error) {
					log.Printf("rpc got result:%v %v\n", result, e)
					if e != nil {
						log.Printf("callback(true)\n")
						callback(true, e)	
					} else {
						callback(false, nil)	
					}
			})))
		}}, func(args ...interface{}) {
			isErr, _ := args[0].(bool)
			log.Printf("args:%+v\n", args)
			log.Printf("isErr:%v\n", isErr)
			if isErr {				
				self.doLoginRet(cbFunc, 1, "error")
			} else {
				self.doLoginRet(cbFunc, 0, "succ")
			}
			
		})
	}()
	
	return nil
}