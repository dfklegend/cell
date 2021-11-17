package service

import (
    "fmt"   
    "strings" 
)


var TheMgr *ServiceMgr

func NewMgr() {
    if TheMgr != nil {
        return
    }
    TheMgr = &ServiceMgr{
        Services: make(map[string]*Service),
        Comps: &Components{},
    }
}

type ServiceMgr struct {  
    Services  map[string]*Service
    Comps *Components
}

func (s *ServiceMgr) Register(c Component, options ...Option) { 
    s.Comps.Register(c, options...)
}

func (s *ServiceMgr) Start() error { 
    // 初始化
    components := s.Comps.List()
    for _, c := range components {
        err := s.newService(c.Comp, c.Opts)
        if err != nil {
            return err
        }
    }

    return nil
}

func (s *ServiceMgr) newService(comp Component, opts []Option) error {
    srv := NewService(comp, opts)

    if _, ok := s.Services[srv.Name]; ok {
        return fmt.Errorf("handler: service already defined: %s", srv.Name)
    }

    if err := srv.ExtractHandler(); err != nil {
        return err
    }

    // 记录
    s.Services[srv.Name] = srv    
    return nil
}

// 1. 转发给对应的Service
// 2. 按目标参数类型反序列化
func (s *ServiceMgr) Call(route string, args []byte) error{
    subs := strings.Split(route, ".")
    if len(subs) != 2 {
        return fmt.Errorf("bat route: %s", route)
    }
    serviceName := subs[0];
    methodName := subs[1];

    service := s.Services[serviceName]
    if service == nil {
        return fmt.Errorf("can not find service: %s", serviceName)
    }
    return service.CallMethod(methodName, args)    
}


// 序列化
func (s *ServiceMgr) AddReqReply(reqId int, args interface{}) {

}
