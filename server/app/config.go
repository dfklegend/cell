package app

import (
    "log"

    "github.com/dfklegend/cell/utils/logger"
    uconfig "github.com/dfklegend/cell/utils/config"
    "github.com/dfklegend/cell/net/server/acceptor"    
    nsession "github.com/dfklegend/cell/net/server/session"
    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/discovery"    
    "github.com/dfklegend/cell/server/config"
    "github.com/dfklegend/cell/server/client/session"
)

type AppConfig struct {
    // 文件配置
    // config.yaml
    cfgFile                    *uconfig.Config
    // servers.yaml
    servers                    *config.ServersCfg
    
	// 自身启动服务器
	Server                     *common.Server
    SessionCfg                 *nsession.SessionConfig
    ETCDConfig                 *discovery.EtcdServiceDiscoveryConfig	  
    WSConfig                   *acceptor.WSConfig    
}

func BuildAppConfig(cfgPath string, serverId string) *AppConfig {
    cfg := &AppConfig{}
    cfg.build(cfgPath, serverId)    
    return cfg
}

// 配置过程
// NewAppConfig
// LoadFileCfgs
// SetupConfig
// SelectServer
func (self *AppConfig) build(cfgPath string, serverId string) {
    self.LoadFileCfgs(cfgPath)
    self.SetupConfig()
    self.SelectServer(serverId)    
}

// 读取文件配置
func (self *AppConfig) LoadFileCfgs(cfgPath string) {
    self.servers = config.LoadServers(cfgPath)
    self.cfgFile = uconfig.LoadConfig(cfgPath)
}

func (self *AppConfig) SetupConfig() {
    // 设置其他子配置配置
    cfgFile := self.cfgFile

    self.SessionCfg = nsession.NewSessionConfig(cfgFile) 
    self.SessionCfg.Impl = session.GetFrontCSImpl() 
      
    self.WSConfig = acceptor.NewWSConfig(cfgFile)
    self.ETCDConfig = discovery.NewEtcdServiceDiscoveryConfig(cfgFile)
}

func (self *AppConfig) SelectServer(serverId string) {
    logger.Log.Infof("start server:%v", serverId)
    one := self.servers.Servers[serverId]
    if one != nil {
        self.Server = one
        return
    }
    logger.Log.Infof("can not find config server:%v", serverId)
    self.Server = &common.Server{}
}

func (self *AppConfig) GetCfgFile() *uconfig.Config {
    return self.cfgFile
}

func (self *AppConfig) GetRPCListenAddr() string {
    return self.Server.Address
}

func (self *AppConfig) GetServerType() string {
    return self.Server.Type
}

func (self *AppConfig) GetServerId() string {
    return self.Server.ID
}

func (self *AppConfig) DumpInfo() {
    log.Printf("Server:%v\n", self.Server)
}