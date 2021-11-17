package client

import (
	"sync"

	"github.com/dfklegend/cell/utils/logger"

	rpccommon "github.com/dfklegend/cell/rpc/common"
)

// 邮局对象，负责向所有已知地址投递消息
// 并获取回信
type MailStation struct {
	// 唯一id 服务器组内唯一
	stationId string
	// serverId:*MailBox
	mailBoxs map[string]*MailBox
	mutex    sync.Mutex
}

func NewMailStation(stationId string) *MailStation {
	return &MailStation{
		stationId: stationId,
		mailBoxs:  make(map[string]*MailBox),
	}
}

func (self *MailStation) GetStationId() string {
	return self.stationId
}

func (self *MailStation) Start() {
}

func (self *MailStation) Stop() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, one := range self.mailBoxs {
		one.Stop()
	}

	self.mailBoxs = make(map[string]*MailBox)
}

// serverId唯一
func (self *MailStation) AddServer(serverId string, address string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if self.mailBoxs[serverId] != nil {
		logger.Log.Infof("duplicated serverId:%v", serverId)
		return
	}

	mb := NewMailBox()
	mb.SetStation(self)
	mb.Start(address)

	self.addServer(serverId, mb)
	logger.Log.Debugf("mailStation addServer:%v %v", serverId, address)
}

func (self *MailStation) addServer(serverId string, mb *MailBox) {
	self.mailBoxs[serverId] = mb
}

func (self *MailStation) RemoveServer(serverId string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	mb := self.mailBoxs[serverId]
	if mb == nil {
		return
	}

	// TODO: 移除服务器前是否把所有pending的rpc返回错误
	mb.Stop()
	self.removeServer(serverId)
	logger.Log.Debugf("mailStation removeServer:%v", serverId)
}

func (self *MailStation) removeServer(serverId string) {
	delete(self.mailBoxs, serverId)
}

func (self *MailStation) getServer(serverId string) *MailBox{
	self.mutex.Lock()
	defer self.mutex.Unlock()
	return self.mailBoxs[serverId]
}

func (self *MailStation) Call(serverId string, method string,
	inArg interface{}, options ...rpccommon.ReqOption) {
	mb := self.getServer(serverId)
	if mb == nil {
		logger.Log.Errorf("can not find server:%v", serverId)
		return
	}

	mb.Call(method, inArg, options...)
}

func (self *MailStation) DumpStat() {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, one := range self.mailBoxs {
		one.DumpStat()
	}
}

// 判断rpc和ack匹配
func (self *MailStation) IsAllBenchOver() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, one := range self.mailBoxs {
		if !one.IsBenchOver() {
			return false
		}
	}
	return true
}
