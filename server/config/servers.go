package config

import (
    "github.com/spf13/viper"

    "github.com/dfklegend/cell/server/common"
)

// 读取servers.yaml
// 对应服务器列表配置
type ServersCfg struct {
    Servers map[string]*common.Server
}

func LoadServers(path string) *ServersCfg {
    v := viper.New()
    v.SetConfigName("servers")
    v.AddConfigPath(path)
    v.SetConfigType("yaml")
    v.ReadInConfig()    

    var obj ServersCfg

    v.Unmarshal(&obj)
    setServersID(&obj)    

    return &obj
}

func setServersID(cfg *ServersCfg) {
    for k, v := range cfg.Servers {
        v.ID = k
    }    
}

