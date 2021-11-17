package channel

import (
    "sync"
    "fmt"

    ucommon "github.com/dfklegend/cell/utils/common"
    api "github.com/dfklegend/cell/rpc/apientry"
    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/interfaces"
)

var (
    channelService = newChannelService()
    tempIdService *ucommon.SerialIdService = ucommon.NewSerialIdService()
)

// 频道服务
type ChannelService struct {
    channels sync.Map
}

func newChannelService() *ChannelService {
    return &ChannelService {        
    }
}

func GetChannelService() *ChannelService {
    return channelService
}

func (self *ChannelService) AddChannel(name string) *Channel {    
    c, ok := self.channels.Load(name)
    if ok {
        return c.(*Channel)
    }
    nc := NewChannel(name)
    self.channels.Store(name, nc)
    return nc
}

func (self *ChannelService) GetChannel(name string) *Channel {    
    c, ok := self.channels.Load(name)
    if !ok {
        return nil
    }
    return c.(*Channel)    
}

func (self *ChannelService) DeleteChannel(name string) {    
    self.channels.Delete(name)
}

func (self *ChannelService) AddToChannel(name string, serverId string,
     netId common.NetIdType) *Channel {    
    c := self.AddChannel(name)
    c.Add(serverId, netId)
    return c
}

func (self *ChannelService) LeaveFromChannel(name string, serverId string,
     netId common.NetIdType) {    
    c := self.GetChannel(name)
    if c != nil {
        c.Leave(serverId, netId)
    }    
}

// 直接向[frontId,netId]推送消息
func (self* ChannelService) PushMessageById(serverId string, netId common.NetIdType,
    route string, msg interface{}, cbFunc api.HandlerCBFunc) {
    sys := interfaces.GetApp().GetSysService()
    sys.PushMessageById(serverId, netId, route, msg, cbFunc)
}

func (self* ChannelService) PushMessageByIds(serverId string, netIds []common.NetIdType,
    route string, msg interface{}, cbFunc api.HandlerCBFunc) {
    sys := interfaces.GetApp().GetSysService()
    sys.PushMessageByIds(serverId, netIds, route, msg, cbFunc)
}

// 可以申请一个临时的channel，用来方便发送请求
func (self *ChannelService) AllocTempChannel() *Channel {    
    name := fmt.Sprintf("_temp_chanel_%v", tempIdService.AllocId())
    return self.AddChannel(name)
}

// c.Add..
// c.PushMessage

func (self *ChannelService) FreeTempChannel(c *Channel) {    
    self.DeleteChannel(c.GetName())
}
