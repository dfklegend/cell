package runservice

import (	
	"reflect"

	"github.com/dfklegend/cell/utils/sche"
	"github.com/dfklegend/cell/utils/timer"
	"github.com/dfklegend/cell/utils/event"
)


// 标准RunService，方便使用
// 拥有timer
// 事件中心
type StandardRunService struct {
	RunService
	TimerMgr 		*timer.TimerMgr
	EventCenter		*event.LocalEventCenter
}

func NewStandardRunService(name string) *StandardRunService {
	return &StandardRunService{
		RunService: *NewRunService(name),
		TimerMgr: timer.NewTimerMgr(),
		EventCenter: event.NewLocalEventCenter(true),
	}
}

func (self *StandardRunService) GetTimerMgr() *timer.TimerMgr {
	return self.TimerMgr
}

func (self *StandardRunService) GetEventCenter() *event.LocalEventCenter {
	return self.EventCenter
}

func (self *StandardRunService) Start() {
	self.RunService.Start()
	// do something else
	self.addTimerSelector()
	self.addEventSelector()
}

func (self *StandardRunService) Stop() {
	// do something else
	self.TimerMgr.Stop()
	self.EventCenter.Clear()

	self.RunService.Stop()
}


func (self *StandardRunService) addTimerSelector() {
	mgr := self.TimerMgr
	selector := self.selector
	selector.AddSelector(sche.NewFuncSelector(reflect.ValueOf(mgr.GetQueue()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			t := v.Interface().(*timer.TimerObj)
			mgr.Do(t)
		}))
}

func (self *StandardRunService) addEventSelector() {
	ec := self.EventCenter
	selector := self.selector
	selector.AddSelector(sche.NewFuncSelector(reflect.ValueOf(ec.GetChanEvent()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			e := v.Interface().(*event.EventObj)
			ec.DoEvent(e)
		}))
}
