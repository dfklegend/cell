package comps

import ( 
    "strings"

    api "github.com/dfklegend/cell/rpc/apientry"
    "github.com/dfklegend/cell/server/component"
    "github.com/dfklegend/cell/server"
    "github.com/dfklegend/cell/utils/logger"

    "github.com/dfklegend/cell/server/examples/chat/services"
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