package channel

import (
    "sync"

    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/interfaces"
)

// [frontId, netId]代表一个连接
type FrontGroup struct {
    NetIds []common.NetIdType
    mutex sync.Mutex
}

func NewFrontGroup() *FrontGroup {
    return &FrontGroup {
        NetIds: make([]common.NetIdType, 0, 128), 
    }
}

func (self* FrontGroup) Add(netId common.NetIdType) {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    self.NetIds = append(self.NetIds, netId)
}

func (self* FrontGroup) FindIndex(netId common.NetIdType) int {   
    for i := 0 ; i < len(self.NetIds); i ++ {
        if netId == self.NetIds[i] {             
            return i
        }
    }
    return -1
}

func (self* FrontGroup) Remove(netId common.NetIdType) {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    i := self.FindIndex(netId)
    // can not find
    if i == -1 {
        return
    }

    if i == 0 {
        self.NetIds = self.NetIds[1:]
        return
    }

    if i == len(self.NetIds)-1 {
        self.NetIds = self.NetIds[:i]
        return
    }
    
    self.NetIds = append(self.NetIds[:i], self.NetIds[i+1:]...)
}

func (self* FrontGroup) Lock() {
    self.mutex.Lock()
}

func (self* FrontGroup) Unlock() {
    self.mutex.Unlock()
}

type Channel struct {
    //groups map[string]*FrontGroup
    //rw sync.RWMutex
    groups sync.Map
    name string
}

func NewChannel(name string) *Channel {
    return &Channel {
        //groups: make(map[string]*FrontGroup),
        name: name,
    }
}

func (self *Channel) GetName() string {
    return self.name
}

func (self *Channel) getGroup(frontId string, createIfMiss bool) *FrontGroup {    
    g, ok := self.groups.Load(frontId)
    if ok {
        return g.(*FrontGroup)
    }
    if !createIfMiss {
        return nil
    }    
    ng := NewFrontGroup()    
    self.groups.Store(frontId, ng)    
    return ng
}

/**
 * 加入频道
 * @param uid{string} 玩家对应的netId
 * @param frontId{string} 前端服务器id
 */
func (self *Channel) Add(frontId string, netId common.NetIdType) {
    g := self.getGroup(frontId, true)
    g.Add(netId)
}

func (self *Channel) Leave(frontId string, netId common.NetIdType) {
    g := self.getGroup(frontId, false)
    if g != nil {
        g.Remove(netId)    
    }    
}

func (self *Channel) Range(f func(k, v interface{}) bool) {
    self.groups.Range(f)
}

func (self *Channel) PushMessage(route string, msg interface{}) {
    sys := interfaces.GetApp().GetSysService()    
    
    self.Range(func( k, v interface{}) bool {
        n := k.(string)
        g := v.(*FrontGroup)

        g.Lock()
        sys.PushMessageByIds(n, g.NetIds, route, msg, nil)
        g.Unlock()
        return true
    })
}

