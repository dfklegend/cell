package scene

import (    
    "fmt"

    "github.com/dfklegend/cell/utils/runservice"
    "github.com/dfklegend/cell/utils/logger" 
)

// run pool
// 
type RunPool struct  {
    name string    
    services []*runservice.StandardRunService
}

func NewRunPool(name string) *RunPool {
    return &RunPool{
        name: name,        
        services: make([]*runservice.StandardRunService, 0),
    }
}

func (self *RunPool) Start(serviceNum int) {    
    for i := len(self.services); i < serviceNum; i ++ {
        name := fmt.Sprintf("%v_%v", self.name, i)
        s := runservice.NewStandardRunService(name)
        s.Start()
        self.services = append(self.services, s)

        logger.Log.Debugf("start scene service:%v", name)
    }
}

func (self *RunPool) GetService(index int) *runservice.StandardRunService {
    if(len(self.services) == 0) {
        return nil
    }
    return self.services[index%len(self.services)]
}