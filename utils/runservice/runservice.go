package runservice

import (
	"log"
	"reflect"
	//"fmt"

	"github.com/dfklegend/cell/utils/sche"
)

// RunService的
var rsScheMgr *sche.ScheMgr

func init() {
	rsScheMgr = sche.DefaultScheMgr
}

func GetScheMgr() *sche.ScheMgr {
	return rsScheMgr
}

// 代表一个可执行的service
// 每个一个单独的routine
// 是一个可投递目标
// 可以用selector来定制循环
type RunService struct {
	Name string
	// 调度器
	scheduler *sche.Sche
	// 用于添加select队列
	selector *sche.MultiSelector
	running  bool

	chanClose chan int
}

func makeScheName(name string) string {
	return name//fmt.Sprintf("rs_%v", name)
}

func NewRunService(name string) *RunService {
	return &RunService{
		Name:      name,
		scheduler: rsScheMgr.GetSche(makeScheName(name)),
		selector:  sche.NewMultiSelector(),
		running:   true,
		chanClose: make(chan int, 1),
	}
}

func (self *RunService) GetScheduler() *sche.Sche {
	return self.scheduler
}

func (self *RunService) GetSelector() *sche.MultiSelector {
	return self.selector
}

func (self *RunService) Start() {
	self.addSchedulerSelector()
	self.addCloseChan()

	go self.loop()
}

func (self *RunService) Stop() {
	// 通知跳出loop
	close(self.chanClose)
	self.scheduler.Stop()
	self.selector.Stop()
	GetScheMgr().DelSche(makeScheName(self.Name))
}

func (self *RunService) addSchedulerSelector() {
	scheduler := self.scheduler
	selector := self.selector
	selector.AddSelector(sche.NewFuncSelector(reflect.ValueOf(scheduler.GetChanTask()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			task := v.Interface().(*sche.RunTask)
			scheduler.DoTask(task)
		}))
}

func (self *RunService) addCloseChan() {
	selector := self.selector
	selector.AddSelector(sche.NewFuncSelector(reflect.ValueOf(self.chanClose),
		func(v reflect.Value, recvOk bool) {
			self.running = false
		}))
}

func (self *RunService) loop() {
	defer func() {
		log.Println("RunServeice loop end")
	}()

	for self.running {
		self.selector.HandleOnce()
	}	
}
