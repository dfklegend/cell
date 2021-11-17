package app

import (
	"reflect"
	"time"

	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/rpc/client"
	"github.com/dfklegend/cell/rpc/server"
	"github.com/dfklegend/cell/utils/logger"
	"github.com/dfklegend/cell/server/component"
	"github.com/dfklegend/cell/net/server/acceptor"
	"github.com/dfklegend/cell/server/discovery"
	pkghandler "github.com/dfklegend/cell/server/handler"	
	ninterfaces "github.com/dfklegend/cell/net/interfaces"
	nsession "github.com/dfklegend/cell/net/server/session"

	"github.com/dfklegend/cell/server/services"
	"github.com/dfklegend/cell/server/interfaces"
	"github.com/dfklegend/cell/server/common"
	"github.com/dfklegend/cell/server/services/sys"
)

type App struct {
	cfg					*AppConfig

	comps 				[]component.IComponent
	acceptors        	[]acceptor.Acceptor

	mailStation 		*client.MailStation
	node        		*server.RPCServerNode

	handler 			*pkghandler.HandlerNode
	msgProcessor 		ninterfaces.IMsgProcessor

	serviceDiscovery    discovery.ServiceDiscovery
	serverListener		*ServerListener

	appDieChan			chan bool
}

/*
	启动流程
		构建配置
		NewApp

		AddComponent
		AddAcceptor

		Start
 */

// 构建好配置后
func NewApp(cfg *AppConfig) *App {
	return &App {
		cfg: cfg,
		comps: make([]component.IComponent, 0),
		acceptors: make([]acceptor.Acceptor, 0),
		mailStation: client.NewMailStation(cfg.Server.ID), 
		node: server.NewNode(cfg.Server.ID),
		handler: pkghandler.NewHandlerNode(), 
		msgProcessor: &services.MsgProcessor{}, 
		appDieChan: make(chan bool),
	}
}

func (self *App) AddComponent(comps ...component.IComponent) {
	for _, c := range comps {
		self.comps = append(self.comps, c)
	}
}

func (self *App) AddAcceptor(acceptors ...acceptor.Acceptor) {
	for _, a := range acceptors {
		self.acceptors = append(self.acceptors, a)
	}
}

func (self *App) GetDieChan() chan bool {
	return self.appDieChan
}

func (self *App) GetCfg() *AppConfig {
	return self.cfg
}

func (self *App) GetMailStation() *client.MailStation {
	return self.mailStation
}

func (self *App) GetHandler() interfaces.IHandlerProcessor {
	return self.handler
}

// func (self *App) GetMsgProcessor() interfaces.IMsgProcessor {
// 	return self.msgProcessor
// }

func (self *App) GetSysService() interfaces.ISysService {
	return sys.GetSysService()
}

func (self *App) GetServerType() string {
	return self.cfg.GetServerType()
}

func (self *App) GetServerId() string {
    return self.cfg.GetServerId()
}

func (self *App) GetServersByType(serverType string) (map[string]*common.Server, error) {
	return self.serviceDiscovery.GetServersByType(serverType)
}

func (self *App) preStart() {
	// add some sys comp
	self.AddComponent(&sys.SysComp{})
	
	nsession.SetMsgProcessor(self.msgProcessor)
}

func (self *App) Start() {	
	self.preStart()

	logger.Log.Debugf("comps:", self.comps)

	for _, c := range self.comps {
		c.Init()
	}

	for _, c := range self.comps {
		c.Start()
	}	
	
	self.mailStation.Start()
	self.node.Start(self.cfg.GetRPCListenAddr())
	self.handler.Start()

	sys.GetSysService().Start()

	// 监听器
	self.listen()
	self.startETCDService()



	self.waitAppEnd()
}

func (self * App) Stop() {
	self.stopETCDService()

	for _, c := range self.comps {
		c.Stop()
	}

	close(self.appDieChan)
}

func (self *App) waitAppEnd() {
	for true {
		select {
		case <-self.appDieChan:
			return
		}
	}
	logger.Log.Infof("app end")
}

func (self *App) startETCDService() {
    cfg := self.GetCfg()
    
    sd, _:= discovery.NewEtcdServiceDiscovery(cfg.ETCDConfig,
        cfg.Server, self.GetDieChan())
    self.serviceDiscovery = sd

    // add listener
    // 如果服务器近期启动，etcd内有残留，先监听
    self.serverListener = MakeServerListener(self, sd)

    sd.Start()    
}

func (self *App) stopETCDService() {
	if self.serviceDiscovery == nil {
		return
	}
	self.serviceDiscovery.Stop()
}

//
func (self *App) RegisterRemoteService(e api.IAPIEntry, options ...api.Option) {
	self.node.GetCollection().Register(e, options...)
}

// 前端接口
func (self *App) RegisterHandlerService(e api.IAPIEntry, options ...api.Option) {
	self.handler.GetCollection().Register(e, options...)
}

//

// server.service.method
// 分析route
// 判断服务类型是否一致
// 一致，调用handler
// 不一致，使用route规则来选择一个服务器
// call rpc
func (self *App) CallHandler(session api.IHandlerSession, route string, msg []byte, cbFunc api.HandlerCBFunc) {
    //
}

func (self *App) listen() {
	// 监听每一个acceptor
	for _, a := range self.acceptors {		
		self.startAcceptorListen(a)
	}
}

func (self *App) startAcceptorListen(a acceptor.Acceptor) {
	// new session
	go func() {
		for conn := range a.GetConnChan() {
			logger.Log.Debugf("new conn come:%v", conn)
			// 新连接建立				
			s := nsession.NewClientSession(conn,
				self.cfg.SessionCfg)
			// add to 
			s.Handle()
			logger.Log.Debugf("%v", conn)
		}
	}()

	go func() {
		// 监听
		a.ListenAndServe()
	}()

	go func() {
		time.Sleep(time.Second)
		logger.Log.Infof("listening with acceptor %s on addr %s", reflect.TypeOf(a), a.GetAddr())	
	}()		
}
