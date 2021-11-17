package cell

import (    
    "github.com/dfklegend/cell/server/app"
    "github.com/dfklegend/cell/net/server/acceptor"    
    "github.com/dfklegend/cell/server/interfaces"
)

// 快捷方式 方便使用
var App *app.App

func createApp(cfg *app.AppConfig) {
    App = app.NewApp(cfg)
    interfaces.SetApp(App)
}

func PrepareApp(cfgPath string, serverId string) {
    cfg := app.BuildAppConfig(cfgPath, serverId)
    createApp(cfg)    

    // cfg is ready

    // add acceptor
    if cfg.Server.Frontend {
        if cfg.Server.WSClientAddress != "" {
            ws := acceptor.NewWSAcceptor(cfg.Server.WSClientAddress,
                cfg.WSConfig.Certs[0], cfg.WSConfig.Certs[1])
            App.AddAcceptor(ws)
        } 
        if cfg.Server.ClientAddress != "" {
            tcp := acceptor.NewTCPAcceptor(cfg.Server.ClientAddress,
                cfg.WSConfig.Certs[0], cfg.WSConfig.Certs[1])
            App.AddAcceptor(tcp)
        }        
    }    
}

// wait setup

func StartApp() {
    App.Start()
}

func StopApp() {    
}