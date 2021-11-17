package main

// 启动routine数量

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/debug"
	rs "github.com/dfklegend/cell/utils/runservice"
	"github.com/dfklegend/cell/utils/sche"
)

func createService(name string) *rs.RunService {
	service := rs.NewRunService(name)
	service.Start()

	return service
}

func createProducer(name string, chanItems chan int) *rs.RunService {
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

func TestServices() {

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

func main() {

	TestServices()
	//debug.DumpCurInfo()
	fmt.Printf("-----------------")
}
