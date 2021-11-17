package sys

import (
    "strings"
    "reflect"

    api "github.com/dfklegend/cell/rpc/apientry"    
    rpccommon "github.com/dfklegend/cell/rpc/common"
    "github.com/dfklegend/cell/utils/logger"
    ucommon "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/utils/runservice"
    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/component"
    "github.com/dfklegend/cell/server/interfaces"
    "github.com/dfklegend/cell/server/client/session"
)

// --------------------
// 协议定义
// 系统请求
type SysReq struct {
    SessionData string `json:"sessiondata"`
    Method string `json:"method"`
    Arg string `json:"arg"`
}

// 返回系统请求
type SysAck struct {
    Result string `json:"result"`
}

type PushSessionReq struct {
    // 推送
    NetId common.NetIdType `json:"netid"`
    SessionData string `json:"sessiondata"`
}

//PushSessionAck SysAck

type KickReq struct {
    // 推送
    NetId common.NetIdType `json:"netid"`    
}

// 向目标推送消息
type PushMessageReq struct {
    NetIds []common.NetIdType `json:"netids"`
    Route string `json:"route"`
    Msg string `json:"msg"`
}

// --------------------


type SysComp struct {	
    component.DefaultComponent
}

func (self* SysComp) Init() {
    app := interfaces.GetApp()

    logger.Log.Debugf("SysComp Init")

    // 内部接口
    app.RegisterRemoteService(&Sys{}, api.WithNameFunc(strings.ToLower),
        api.WithSchedulerName("sysservice"))
}


// --------------------
type Sys struct {
    api.APIEntry
}

func (self *Sys) Call(msg *SysReq, cbFunc api.HandlerCBFunc) error {    
    return GetSysService().Call(msg, cbFunc)
}

// 更新session数据
func (self *Sys) PushSession(msg *PushSessionReq, cbFunc api.HandlerCBFunc) error {    
    return GetSysService().PushSession(msg, cbFunc)
}

func (self *Sys) Kick(msg *KickReq, cbFunc api.HandlerCBFunc) error {    
    return GetSysService().Kick(msg, cbFunc)
}

func (self *Sys) PushMessage(msg *PushMessageReq, cbFunc api.HandlerCBFunc) error {    
    return GetSysService().OnPushMessage(msg, cbFunc)
}


// --------------------
var sysService *SysService = NewSysService()
type SysService struct {    
    runService *runservice.RunService
}

func NewSysService() *SysService {
    service := &SysService {
        runService: runservice.NewRunService("sysservice"),
    }
    return service
}

func GetSysService() *SysService {
    return sysService
}

func (self* SysService) Start() {
    self.runService.Start()
}

// 处理handler接口的调用
func (self* SysService) Call(msg *SysReq, cbFunc api.HandlerCBFunc) error {
    // . 构建backSession对象
    // . 调用handler
    // . 返回值    
    logger.Log.Debugf("msg:%+v", msg)

    // TODO: pool
    // session禁止保存作为后续用途
    // 可以保存[frontId,netId]
    session := session.NewBackSession(msg.SessionData)
    logger.Log.Debugf("backsession:%+v", session)
    
    interfaces.GetApp().GetHandler().Call(session,
        msg.Method, []byte(msg.Arg), func(e error, result interface{}) {
        logger.Log.Debugf("handler.call:%v cb result:%v", msg.Method, result)
        
        if e != nil {
            api.CheckInvokeCBFunc(cbFunc, e,
                &SysAck{                                
            })      
        } else {
            data := result.([]byte)
            api.CheckInvokeCBFunc(cbFunc, nil,
                &SysAck{            
                    Result: string(data),
            })      
        }

        
    })
    return nil
}

// 推送到对应的前端服务器
func (self* SysService) ReqPushSession(serverId string, netId common.NetIdType,
    sessionData string, cbFunc api.HandlerCBFunc) {
    // 调用mailstation
    msg := &PushSessionReq{
        NetId: netId,
        SessionData: sessionData,
    }

    ms := interfaces.GetApp().GetMailStation()
    if cbFunc != nil {
        ms.Call(serverId, "sys.pushsession", msg, rpccommon.ReqWithCB(
            reflect.ValueOf(func(result *SysAck, e error){
                api.CheckInvokeCBFunc(cbFunc, nil, nil) 
            }) ))
    } else {
        ms.Call(serverId, "sys.pushsession", msg)
    }    
}

func (self* SysService) PushSession(msg *PushSessionReq, cbFunc api.HandlerCBFunc) error {    
    logger.Log.Debugf("sys.PushSession:%+v", msg)
    fs := session.GetFrontSessions().FindSession(msg.NetId)
    if fs == nil {
        logger.Log.Errorf("can not find frontsession:%vn", msg.NetId)
        api.CheckInvokeCBFunc(cbFunc, nil, nil) 
        return nil
    }
    fs.Lock()
    fs.FromJson(msg.SessionData)
    fs.Unlock()
    api.CheckInvokeCBFunc(cbFunc, nil, nil) 

    logger.Log.Debugf("fs:%+v fs.Data:%+v\n", fs, fs.Data)
    return nil
}

func (self* SysService) ReqKick(serverId string, netId common.NetIdType,
    cbFunc api.HandlerCBFunc) {
    msg := &KickReq{
        NetId: netId,
    }
    ms := interfaces.GetApp().GetMailStation()
    if cbFunc != nil {
        ms.Call(serverId, "sys.kick", msg, rpccommon.ReqWithCB(
            reflect.ValueOf(func(result *SysAck, e error){
                api.CheckInvokeCBFunc(cbFunc, nil, nil) 
            }) ))
    } else {
        ms.Call(serverId, "sys.kick", msg)
    }    
}

func (self* SysService) Kick(msg *KickReq, cbFunc api.HandlerCBFunc) error {    
    logger.Log.Debugf("sys.Kick:%+v", msg)
    fs := session.GetFrontSessions().FindSession(msg.NetId)
    if fs == nil {
        logger.Log.Errorf("can not find frontsession:%vn", msg.NetId)
        api.CheckInvokeCBFunc(cbFunc, nil, nil) 
        return nil
    }
    fs.Kick()
    api.CheckInvokeCBFunc(cbFunc, nil, nil)     
    return nil
}

func (self* SysService) PushMessageById(serverId string, netId common.NetIdType,
    route string, msg interface{}, cbFunc api.HandlerCBFunc) {
    ids := make([]common.NetIdType, 0)
    ids = append(ids, netId)
    self.PushMessageByIds(serverId, ids, route, msg, cbFunc)
}

func (self* SysService) PushMessageByIds(serverId string, netIds []common.NetIdType,
    route string, msg interface{}, cbFunc api.HandlerCBFunc) {

    // msg判断不是string，用json序列化
    str, ok := msg.(string)
    if !ok {
        str = ucommon.SafeJsonMarshal(msg)
    }
    
    req := &PushMessageReq{
        NetIds: netIds,
        Route: route,
        Msg: str,
    }

    app := interfaces.GetApp()
    // 本服
    if serverId == app.GetServerId() {
        self.OnPushMessage(req, nil)
        api.CheckInvokeCBFunc(cbFunc, nil, nil)
        return
    }

    ms := interfaces.GetApp().GetMailStation()
    if cbFunc != nil {
        ms.Call(serverId, "sys.pushmessage", req, rpccommon.ReqWithCB(
            reflect.ValueOf(func(result *SysAck, e error){
                api.CheckInvokeCBFunc(cbFunc, nil, nil) 
            }) ))
    } else {
        ms.Call(serverId, "sys.pushmessage", req)
    }    
}

func (self* SysService) OnPushMessage(msg *PushMessageReq, cbFunc api.HandlerCBFunc) error {    
    logger.Log.Debugf("sys.OnPushMessage:%+v", msg)

    // 向每一个玩家推送
    for _, v := range(msg.NetIds) {
        fs := session.GetFrontSessions().FindSession(v)
        if fs == nil {
            continue
        }
        fs.Session.Push(msg.Route, []byte(msg.Msg))
    }

    api.CheckInvokeCBFunc(cbFunc, nil, nil)     
    return nil
}