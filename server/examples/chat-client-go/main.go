package main

import (    
    "flag"   
    "fmt"
    "os"
    "os/signal"    

    "github.com/dfklegend/cell/utils/logger"
    "github.com/dfklegend/cell/utils/cmd"
    "github.com/dfklegend/cell/utils/common"

    "github.com/dfklegend/cell/server/examples/chat-client-go/client"
)

func StartCellClient() {
    c := client.NewChatClient("client")
    c.Start("127.0.0.1:30021")    
}

func main() {
    var serverId = flag.String("id", "gate-1", "the server id")
    flag.Parse()
    
    logger.SetDebugLevel()
    logger.Log.Debugf("s:%v", *serverId)
    cmd.StartConsoleCmd()

    StartCellClient()
    common.GoPprofServe("6060")

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, os.Kill)

    s := <-c
    fmt.Println("Got signal:", s)
}
