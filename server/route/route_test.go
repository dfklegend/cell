package route

import (
    "fmt"
    "log"
    "testing"
    //"reflect"

    "github.com/dfklegend/cell/utils/logger"
    sessioni "github.com/dfklegend/cell/server/client/session"
)


func TestErrorParam(t *testing.T) {
    fmt.Println("TestErrorParam") 
    log.Printf("route result:%v\n", TheRouteService.Route("test", 1))
    fmt.Println("--") 
}

func TestMapParam(t *testing.T) {
    s := TheRouteService    

    fmt.Println("TestMapParam") 
    m := make(map[string]interface{})
    m["haha"] = 1       
    log.Printf("route result:%v\n", TheRouteService.Route("test", m))

    s.Register("test", RouteF)
    log.Printf("route result:%v\n", TheRouteService.Route("test", m))    
    fmt.Println("--") 
}

func TestSessionParam(t *testing.T) {
    s := TheRouteService    

    fmt.Println("TestSessionParam") 
    session := sessioni.NewBackSessionForPush("s1", 1)

    //si := session.(sessioni.IServerSession)
    // TypeOfSession := reflect.TypeOf((*sessioni.IServerSession)(nil)).Elem()
    // log.Printf("%v", reflect.TypeOf(session).Implements(TypeOfSession))    

    log.Printf("route result:%v\n", TheRouteService.Route("test", session))

    s.Register("test", RouteF)
    session.Set("id", 1)
    log.Printf("route result:%v\n", TheRouteService.Route("test", session))    
    fmt.Println("--") 
}

func RouteF(serverType string, p IRouteParam) string {
    var id = p.Get("id", 0).(int)
    servers, err := GetServers(serverType)

    if servers == nil || len(servers) == 0 || err != nil {
        logger.Log.Errorf("can not find serverType:%v\n", serverType)
        return "null"
    }

    want := fmt.Sprintf("s%v", id)

    for k, _ := range(servers) {
        if want == k {
            return want
        }
    }

    log.Printf("can not find server:%v\n", want)
    return ""
}

func TestRoute(t *testing.T) {
    fmt.Println("TestRoute") 
    s := TheRouteService

    s.Register("test", RouteF)
    m := make(map[string]interface{})
    m["id"] = 1       
    log.Printf("route result:%v\n", s.Route("test", m))

    m["id"] = 2       
    log.Printf("route result:%v\n", s.Route("test", m))

    m["id"] = 3       
    log.Printf("route result:%v\n", s.Route("test", m))
    fmt.Println("--") 
}

