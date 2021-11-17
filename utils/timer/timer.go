package timer

import (	
	"time"
	"sync"
	//"log"

	"github.com/dfklegend/cell/utils/common"
)

/**
 * timer到期后，被push到channel 
 * 便于整合到具体使用线程
 */

type TimerIdType uint64 
type CBFunc func(...interface{})
type QueueChan chan *TimerObj

var timerMgr = NewTimerMgr()

type ITimerMgr interface {
	After(duration time.Duration, cb CBFunc, args ...interface{}) TimerIdType
	AddTimer(duration time.Duration, cb func(), args ...interface{}) TimerIdType
	Cancel(id TimerIdType)
	GetQueue() QueueChan
}

type TimerObj struct {
	TimerId 	TimerIdType
	Duration 	time.Duration
	CB 			CBFunc
	Args 		[]interface{}
	Canceled 	bool
	timer    	*time.Timer
}

// TODO: 可以考虑pool
func NewTimerObj(id TimerIdType,
	duration time.Duration,
	cb CBFunc,
	args []interface{}) *TimerObj {
	return &TimerObj {
		TimerId: id,
		Duration: duration,
		CB: cb,
		Args: args,
		Canceled: false,
	}
}

// 每个环境可以创建自己的timerMgr
type TimerMgr struct {
	idService 		*common.SerialIdService64
	queue 			QueueChan

	// 用于cancel
	// TimerIdType:*TimerObj
	timers			sync.Map
	running 		bool
}

func NewTimerMgr() *TimerMgr {
	return &TimerMgr {
		idService: common.NewSerialIdService64(),
		queue: make(QueueChan, 999),		
		running: true,
	}
}

// 用于简单测试
// 用户可以创建自己的
func GetTimerMgr() *TimerMgr {
	return timerMgr
}

func (self *TimerMgr) Stop() {
	self.running = false
}

func (self *TimerMgr) allocId() TimerIdType {
	return TimerIdType(self.idService.AllocId())
}

func (self *TimerMgr) After(duration time.Duration, cb CBFunc, args ...interface{}) TimerIdType {
	t := NewTimerObj(self.allocId(),
		0, cb, args)

	self.doLater(duration, t)
	self.timers.Store(t.TimerId, t)	
	return t.TimerId
}

func (self *TimerMgr) AddTimer(duration time.Duration, cb CBFunc, args ...interface{}) TimerIdType {
	t := NewTimerObj(self.allocId(),
		duration, cb, args)

	self.doLater(duration, t)
	self.timers.Store(t.TimerId, t)	
	return t.TimerId
}

func (self *TimerMgr) Cancel(id TimerIdType) {
	v, ok := self.timers.Load(id)
	if !ok {
		return
	}	
	t := v.(*TimerObj)
	t.Canceled = true
	if t.timer != nil {
		t.timer.Stop()
	}

	self.timers.Delete(id)
}

func (self *TimerMgr) doLater(duration time.Duration, t *TimerObj) {
	t.timer = time.AfterFunc(duration, func() {
		//log.Println("cb")
		if t.Canceled {
			return
		}

		// 停止了
		if !self.running {
			return
		}
		// 底层在新routine执行
		// time.AfterFunc
		self.queue <- t
	})
}

/*
	select {
		case t :<- mgr.GetQueue():
			mgr.Do(t)
	}
 */
// 外部从queue读取后，调用
func (self *TimerMgr) Do(t *TimerObj) {
	if t.Canceled {
		return
	}	
	
	t.CB(t.Args...)
	// CB内被cancel掉了
	if t.Canceled {
		return
	}	

	if t.Duration > 0 {
		self.doLater(t.Duration, t)
	} else {
		self.timers.Delete(t.TimerId)
	}
}

func (self *TimerMgr) GetQueue() QueueChan {
	return self.queue
}