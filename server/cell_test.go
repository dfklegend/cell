package cell

import (    
    "log"
    "testing"    
)

func startApp(serverId string) {
    PrepareApp("./data/config", serverId)   
    log.Printf("--\n") 
    App.GetCfg().DumpInfo()
    StartApp()
}

func TestStartApp1(t *testing.T) {
    startApp("logic-1")
}

func _TestStartApp2(t *testing.T) {
    startApp("logic-2")
}


