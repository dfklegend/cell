package main

import (    
    "log"
    "flag"
    //"strings"
    "fmt"

    "github.com/dfklegend/cell/utils/cmd"
    "github.com/dfklegend/cell/utils/logger"
    "github.com/dfklegend/cell/server"    
    "github.com/dfklegend/cell/server/route" 
    "github.com/dfklegend/cell/server/examples/chat/comps"   
)

func RouteChat(serverType string, p route.IRouteParam) string {
    var id = p.Get("chatid", 0).(int)
    servers, err := route.GetServers(serverType)

    if servers == nil || len(servers) == 0 || err != nil {
        logger.Log.Errorf("can not find serverType:%v\n", serverType)
        return ""
    }

    want := fmt.Sprintf("chat-%v", id)

    for k, _ := range(servers) {
        if want == k {
            return want
        }
    }

    log.Printf("can not find server:%v\n", want)
    return ""
}

func startApp(serverId string) {    
    cell.PrepareApp("./data/config", serverId)   
    app := cell.App

    log.Printf("--\n") 
    app.GetCfg().DumpInfo()

    // add component
    serverType := app.GetServerType()
    isChat := serverType == "chat"
    if isChat {
        app.AddComponent(&comps.CompChat{})    
    } else {
        app.AddComponent(&comps.CompGate{})    
    }
    
    route.GetRouteService().Register("chat", RouteChat)

    cell.StartApp()    
}

// 启动gate和chat至少各一
// chat.ext -id=chat-1
// 服务器id参看 data/config/servers.yaml
func main() {
    var serverId = flag.String("id", "gate-1", "the server id")
    flag.Parse()
    
    cmd.StartConsoleCmd()
    startApp(*serverId)
}
