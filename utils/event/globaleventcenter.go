package event

import (
	"sync"
	"sync/atomic"

	//"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/logger"
)

var globalEventCenter = newGlobalEC()

type LocalECList struct {	
	localECs 		sync.Map
	Size			int32
}

func newLocalECList() *LocalECList {
	return & LocalECList{}
}

func (self *LocalECList) Add(ec ILocalEventCenter) {
	self.localECs.Store(ec.GetId(), ec)
	atomic.AddInt32(&self.Size, 1)
}

func (self *LocalECList) Del(ec ILocalEventCenter) {
	self.localECs.Delete(ec.GetId())
	atomic.AddInt32(&self.Size, -1)
}

func (self *LocalECList) Range(f func(k, v interface{}) bool) {
	self.localECs.Range(f)
}


type GlobalEventCenter struct {
	events 			sync.Map	
}

func newGlobalEC() *GlobalEventCenter {
	return & GlobalEventCenter{		
	}
}

func GetGlobalEC() *GlobalEventCenter {
	return globalEventCenter
}

func (self *GlobalEventCenter) getECList(eventName string, createIfMiss bool) *LocalECList {
	v, ok := self.events.Load(eventName)
	var l *LocalECList
	if !ok {
		if !createIfMiss {
			return nil
		}
		l = newLocalECList()
		self.events.Store(eventName, l)
		return l
	}
	return v.(*LocalECList)
}


func (self *GlobalEventCenter) Subscribe(eventName string, child ILocalEventCenter) {
	listeners := self.getECList(eventName, true)	
	listeners.Add(child)
}

func (self *GlobalEventCenter) Unsubscribe(eventName string, child ILocalEventCenter) {
	listeners := self.getECList(eventName, false)
	if listeners == nil {
		return
	}	
	listeners.Del(child)
}

func (self *GlobalEventCenter) Publish(eventName string, args ...interface{}) {
	listeners := self.getECList(eventName, false)
	if listeners == nil {
		return
	}

	e := &EventObj{
		EventName: eventName,
		Args: args,
	}
	listeners.Range(func(k, v interface{}) bool {
		l := v.(ILocalEventCenter)

		select {
		case l.GetChanEvent() <- e:
			return true
		default:
			// 警告，某个事件队列一直没收取消息
			logger.Log.Warnf("localEventCenter:%v queue full", l.GetId())
			return true
		}
		return true
	})
}

func (self *GlobalEventCenter) DumpInfo() {
	logger.Log.Debugf("GlobalEventCenter info:")
	self.events.Range(func(k, v interface{}) bool {
		en := k.(string)
		el := v.(*LocalECList)
		logger.Log.Debugf("%v size:%v", en, el.Size)
		return true
	})
}