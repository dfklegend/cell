package event

import (	
	"sync"

	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/logger"
)

var (
	idService = common.NewSerialIdService64()
	ecIdService = common.NewSerialIdService64()
	useFakeMutex = false
)

func NewRWMutex() common.IMutex {	
	if useFakeMutex {
		return &common.FakeMutex{}	
	}	
	return &sync.RWMutex{}
}

// --------

type ListenerList struct {	
	Global 				bool
	RWLock 				common.IMutex
	Listeners 			map[uint64]*EventListener	
}

func NewListenerList() *ListenerList {
	return &ListenerList{
		RWLock: NewRWMutex(),
		Listeners: make(map[uint64]*EventListener),
	}
}

func (self *ListenerList) Add(l *EventListener) {
	self.RWLock.Lock()
	defer self.RWLock.Unlock()
	self.Listeners[l.Id] = l
}

func (self *ListenerList) Del(id uint64) {
	self.RWLock.Lock()
	defer self.RWLock.Unlock()
	delete(self.Listeners, id)
}

func (self *ListenerList) GetSize() int {
	return len(self.Listeners)
}

// --------

func NewEventListener() *EventListener {
	return &EventListener{
		Id: idService.AllocId(),
		Args: make([]interface{}, 0),		
	}
}

type LocalEventCenter struct {
	id 			uint64
	RWLock 		common.IMutex
	events 		map[string]*ListenerList
	// 接收global
	chanEvent   ChanEvent
	// 是否都使用chan
	// 包括本地事件
	localUseChan bool
	running bool
}

func NewLocalEventCenter(useChan bool) *LocalEventCenter {
	return &LocalEventCenter {
		id: ecIdService.AllocId(),
		events: make(map[string]*ListenerList),
		chanEvent: make(ChanEvent, 999),
		localUseChan: useChan,
		RWLock: NewRWMutex(), 
		running: true,
	}
}

func (self *LocalEventCenter) Clear() {
	// 需要取消所有全局的事件注册
	self.running = false
	self.RWLock.Lock()
	defer self.RWLock.Unlock()

	for en, el := range(self.events) {
		if !el.Global {
			continue
		}

		GetGlobalEC().Unsubscribe(en, self)
	}

	self.events = make(map[string]*ListenerList)
}

func (self *LocalEventCenter) SetLocalUseChan(v bool) {
	self.localUseChan = v
}

func (self *LocalEventCenter) GetId() uint64 {
	return self.id
}

func (self *LocalEventCenter) getListener(eventName string, createIfMiss bool) *ListenerList {
	self.RWLock.RLock()
	v, ok := self.events[eventName]
	self.RWLock.RUnlock()
	if !ok {
		if !createIfMiss {
			return nil
		}
		v = NewListenerList()

		self.RWLock.Lock()
		if self.running {
			self.events[eventName] = v	
		}		
		self.RWLock.Unlock()
	}
	return v
}

func (self *LocalEventCenter) Subscribe(eventName string, cb CBFunc, args ...interface{}) uint64 {
	if !self.running {
		return 0
	}
	list := self.getListener(eventName, true)
	l := NewEventListener()	
	l.Args = args
	l.CB = cb
	list.Add(l)
	return l.Id
}

// 先注册自己
// 向Global注册
func (self *LocalEventCenter) GSubscribe(eventName string, cb CBFunc, args ...interface{}) uint64 {
	if !self.running {
		return 0
	}
	list := self.getListener(eventName, true)
	id := self.Subscribe(eventName, cb, args...)
	if !list.Global {
		// 向global注册
		GetGlobalEC().Subscribe(eventName, self)
		list.Global = true
	}
	return id 
}

func (self *LocalEventCenter) Unsubscribe(eventName string, id uint64) {
	list := self.getListener(eventName, false)
	if list == nil {
		return
	}
	list.Del(id)

	// 发现自己空了，global节点取消
	if list.Global && list.GetSize() == 0 {
		GetGlobalEC().Unsubscribe(eventName, self)
		list.Global = false
	}
}

func (self *LocalEventCenter) GUnsubscribe(eventName string, id uint64) {
	self.Unsubscribe(eventName, id)
}

func (self *LocalEventCenter) GetChanEvent() ChanEvent {
	return self.chanEvent
}

// call by chan receiver
func (self *LocalEventCenter) DoEvent(e *EventObj) {
	// 派发
	self.dispatch(e.EventName, e.Args...)
}

func (self *LocalEventCenter) Publish(eventName string, args ...interface{}) {
	if self.localUseChan {
		e := &EventObj{
			EventName: eventName,
			Args: args,
		}
		self.chanEvent <- e
	} else {
		self.dispatch(eventName, args...)
	}
}

func (self *LocalEventCenter) dispatch(eventName string, args ...interface{}) {
	ll := self.getListener(eventName, false)
	if ll == nil {
		return
	}
	listeners := ll.Listeners

	ll.RWLock.RLock()
	for _, v := range(listeners) {
		finalArgs := append(v.Args, args...)
		v.CB(finalArgs...)
	}
	ll.RWLock.RUnlock()
}

func (self *LocalEventCenter) DumpInfo(eventName string) {
	ll := self.getListener(eventName, false)
	if ll == nil {
		return
	}

	ll.RWLock.RLock()
	logger.Log.Debugf("%v size:%v", eventName, ll.GetSize())
	ll.RWLock.RUnlock()
}