package app

import (    
    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/discovery"
    "github.com/dfklegend/cell/utils/logger"
)

// 监听服务发现

type ServerListener struct {
    app                 *App    
}

func MakeServerListener(app *App,
    serviceDiscovery discovery.ServiceDiscovery) *ServerListener {
    listener := &ServerListener {
        app: app, 
    }

    serviceDiscovery.AddListener(listener)
    return listener
}

func (self *ServerListener) AddServer(server *common.Server) {
    logger.Log.Debugf("ServerListener add server:%v", server)
    // add to mailstation
    mb := self.app.GetMailStation()
    mb.AddServer(server.ID, server.Address)
}

func (self *ServerListener) RemoveServer(server *common.Server) {
    logger.Log.Debugf("ServerListener remove server:%v", server)
    mb := self.app.GetMailStation()
    mb.RemoveServer(server.ID)
}

