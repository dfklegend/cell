package main

import (    
    "log"
    "flag"
    "net/http"
    _ "net/http/pprof"

    "github.com/dfklegend/cell/utils/cmd"
    "github.com/dfklegend/cell/server"

    "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene"
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space/validators"
)

func startApp(serverId string) {

    cell.PrepareApp("./data/config", serverId)   
    app := cell.App

    log.Printf("--\n") 
    app.GetCfg().DumpInfo()

    // add component
    //serverType := app.GetServerType()
    //isChat := serverType == "chat"  
    
    scene.GetSceneMgr().Start()
    cell.StartApp()        
}

func pprofServe() {
    http.ListenAndServe("0.0.0.0:6060", nil)
}

func main() {
    validators.Visit()

    var serverId = flag.String("id", "gate-1", "the server id")    
    var oneRoutine = flag.Bool("one", true, "all scene on routine")
    flag.Parse()
    
    // 单独routine
    scene.SetAllSceneOnRoutine(*oneRoutine)
    cmd.StartConsoleCmd()
    go pprofServe()
    startApp(*serverId)
}
