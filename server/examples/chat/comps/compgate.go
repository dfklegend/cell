package comps

import ( 
    "strings"

    api "github.com/dfklegend/cell/rpc/apientry"
    "github.com/dfklegend/cell/server/component"
    "github.com/dfklegend/cell/server"
    "github.com/dfklegend/cell/utils/logger"

    "github.com/dfklegend/cell/server/examples/chat/services"
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