package runservice

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/debug"
	"github.com/dfklegend/cell/utils/sche"
)

func DoPost(service *RunService, index int) {
	fmt.Printf("post:%v\n", index)
	service.GetScheduler().Post(func() (interface{}, error) {
		fmt.Printf("run:%v\n", index)		
		return nil, nil
	})
}

// 模拟启动一个RunService
func Test1_RunService(t *testing.T) {
	service := NewRunService("test")
	service.Start()

	running := true
	go func() {
		index := 1
		for running {
			DoPost(service, index)
			index ++
		}

	}()

	go func() {
		time.Sleep(1 * time.Second)
		service.Stop()
	}()

	time.Sleep(2 * time.Second)
	running = false
	debug.DumpCurInfo()
	fmt.Println("---- Test_RunService over ----")
}

func createService(name string) *RunService {
	service := NewRunService(name)
	service.Start()

	return service
}

func createProducer(name string, chanItems chan int) *RunService {
	service := createService(name)
	// 创建一个timer
	t := time.NewTicker(1 * time.Second)

	service.GetSelector().AddSelector(
		sche.NewFuncSelector(reflect.ValueOf(t.C),
			func(v reflect.Value, recvOk bool) {
				// 生产一个数据
				chanItems <- 1
				fmt.Printf("%v %v:%v\n", common.GetRoutineID(), name, 1)
			}))

	go func() {
		time.Sleep(3 * time.Second)
		service.Stop()
		t.Stop()
		fmt.Printf("%v %v stop\n", common.GetRoutineID(), name)
	}()
	return service
}

func createConsumer(name string, chanItems chan int) {
	service := createService(name)

	service.GetSelector().AddSelector(
		sche.NewFuncSelector(reflect.ValueOf(chanItems),
			func(v reflect.Value, recvOk bool) {
				i := v.Int()
				fmt.Printf("%v %v:%v\n", common.GetRoutineID(), name, i)
				time.Sleep(10 * time.Millisecond)
			}))

	go func() {
		time.Sleep(3 * time.Second)
		service.Stop()
		fmt.Printf("%v %v stop\n", common.GetRoutineID(), name)
	}()
}

func Test_Services(t *testing.T) {

	chanItems := make(chan int, 999999)

	// 生产者每秒产生一个
	for i := 0; i < 2; i++ {
		createProducer(fmt.Sprintf("producer%v", i), chanItems)
	}

	for i := 0; i < 1; i++ {
		createConsumer(fmt.Sprintf("consumer%v", i), chanItems)
	}

	// go func() {
	// 	time.Sleep(1 * time.Second)
	// }()

	time.Sleep(5 * time.Second)
	debug.DumpCurInfo()
	fmt.Println("---- Test_Services over ----")
}

func Test_RoutineNum(t *testing.T) {
	debug.DumpCurInfo()
	fmt.Println("---- Test_RoutineNum over ----")
}

func Test_Standard(t *testing.T) {
	log.Printf("---- Test_Standard ----")	
	service := NewStandardRunService("test")
	service.Start()

	total := 0
	t1R := 0
	service.GetTimerMgr().AddTimer(100*time.Millisecond, func(args ...interface{}) {
		t1R = common.GetRoutineID()
		log.Printf("t2 r:%v %v\n", t1R, args)
		
		total += 2
	}, 3, 5)

	center := service.GetEventCenter()
	service.GetEventCenter().Subscribe("event1", func(args ...interface{}) {
		log.Printf("on event1")	
	})

	go func(){
		time.Sleep(1 * time.Second)
		center.Publish("event1",1)		
	}()

	time.Sleep(3 * time.Second)
	service.Stop()
	time.Sleep(3 * time.Second)

	log.Printf("over")	
}
