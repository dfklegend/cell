package comps

import ( 
    "strings"

    api "dfk.com/cell/rpc/apientry"
    "dfk.com/cell/server/component"
    "dfk.com/cell/server"
    "dfk.com/cell/utils/logger"

    "dfk.com/cell/server/examples/chat/services"
)

type CompChat struct {
    component.DefaultComponent
}


func (self *CompChat) Init() {
    // 注册各种service
    app := cell.App

    logger.Log.Debugf("CompChat Init")
    app.RegisterHandlerService(&services.Chat{}, api.WithNameFunc(strings.ToLower),
        api.WithSchedulerName("chatservice"))
    app.RegisterRemoteService(&services.ChatRemote{}, api.WithNameFunc(strings.ToLower),
        api.WithSchedulerName("chatservice"))
}

func (self *CompChat) Start() {
    services.GetChatService().Start()
}