package common

import (
	"log"
	"sync"	
)

type WaitAckReq struct {
	ServerReqId ReqIdType
	ClientId    string
	ClientReqId ReqIdType
	// TODO:时间戳 过期移除
	// 是否获得ack
	GotAck  bool
	Err     string
	Message string
}

func NewWaitAckReq() *WaitAckReq {
	return &WaitAckReq{
		GotAck: false,
	}
}

// ----------------------------

// 管理器
// TODO: 超时移除
type WaitAckCtrl struct {
	mutex sync.Mutex
	acks  map[ReqIdType]*WaitAckReq
}

func NewWaitAckCtrl() *WaitAckCtrl {
	return &WaitAckCtrl{
		acks: make(map[ReqIdType]*WaitAckReq),
	}
}

// 实际的
func (self *WaitAckCtrl) RegisterWaitAck(wa *WaitAckReq) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.acks[wa.ServerReqId] = wa

	//log.Printf("RegisterWaitAck:%v\n", wa)
}

func (self *WaitAckCtrl) RemoveWaitAck(serverReqId ReqIdType) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	wa := self.acks[serverReqId]
	if wa == nil {
		return
	}
	delete(self.acks, serverReqId)
}

func (self *WaitAckCtrl) PushAck(serverReqId ReqIdType, err string, message string) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	wa := self.acks[serverReqId]
	if wa == nil {
		log.Printf("error can not find waitack:%v\n", serverReqId)
		return
	}

	wa.GotAck = true
	wa.Err = err
	wa.Message = message
}

// pop 获得返回的acks
func (self *WaitAckCtrl) PopReadyAcks(clientId string, maxNum int) []*WaitAckReq {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if len(self.acks) == 0 {
		return nil
	}

	rets := make([]*WaitAckReq, 0)
	num := 0
	for _, wa := range self.acks {
		if num >= maxNum {
			break
		}
		if clientId == wa.ClientId && wa.GotAck {
			rets = append(rets, wa)
			num++
		}
	}

	for _, wa := range rets {
		delete(self.acks, wa.ServerReqId)
	}
	//log.Printf("pop acks:%v\n", rets)
	return rets
}
