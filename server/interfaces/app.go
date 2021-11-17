package interfaces

import (
    "github.com/dfklegend/cell/rpc/client"
    api "github.com/dfklegend/cell/rpc/apientry"
    "github.com/dfklegend/cell/server/component"
    "github.com/dfklegend/cell/server/common"
    //"github.com/dfklegend/cell/net/common/conn/message"
)

// 定义方便引用
// type IMsgProcessor interface{    
//     ProcessMessage(api.IHandlerSession, *message.Message)
// }

type IHandlerProcessor interface {
    Call(session api.IHandlerSession, method string, msg []byte, cbFunc api.HandlerCBFunc)
}

type ISysService interface {
    ReqPushSession(serverId string, netId common.NetIdType, sessionData string, cbFunc api.HandlerCBFunc)
    ReqKick(serverId string, netId common.NetIdType, cbFunc api.HandlerCBFunc)
    // msg string 或者 结构(将自动json序列化)
    PushMessageById(serverId string, netId common.NetIdType,
        route string, msg interface{}, cbFunc api.HandlerCBFunc)
    PushMessageByIds(serverId string, netIds []common.NetIdType,
        route string, msg interface{}, cbFunc api.HandlerCBFunc)
}

type IApp interface {
    // 启动
    // 添加component
    AddComponent(comp ...component.IComponent)

    Start()
    Stop()

    GetMailStation() *client.MailStation
    //GetMsgProcessor() IMsgProcessor   
    GetHandler() IHandlerProcessor
    GetSysService() ISysService
    GetServerId() string
    GetServerType() string
    GetServersByType(serverType string) (map[string]*common.Server, error)

    RegisterRemoteService(e api.IAPIEntry, options ...api.Option)
    RegisterHandlerService(e api.IAPIEntry, options ...api.Option)
}