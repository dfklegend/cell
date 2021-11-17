package comps

import ( 
    "strings"

    api "dfk.com/cell/rpc/apientry"
    "dfk.com/cell/server/component"
    "dfk.com/cell/server"
    "dfk.com/cell/utils/logger"

    "dfk.com/cell/server/examples/chat/services"
)

type CompGate struct {
    component.DefaultComponent
}


func (self *CompGate) Init() {
    // 注册各种service
    app := cell.App

    logger.Log.Debugf("CompGate Init")
    app.RegisterHandlerService(&services.Gate{}, api.WithNameFunc(strings.ToLower))
}

func (self *CompGate) Start() {    
}