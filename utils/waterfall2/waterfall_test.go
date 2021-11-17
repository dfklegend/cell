package waterfall

import (
	"fmt"
	"testing"
	"time"

	"github.com/dfklegend/cell/utils/common"
	scheP "github.com/dfklegend/cell/utils/sche"
)

func Test_Simple(t *testing.T) {
	fmt.Println("---- show waterfall callback ----")
	sche := scheP.NewSche()
	fn := func(tasks []Task, final Task) {
		Waterfall(sche, tasks, final)
	}
	go test_do(fn, 100)
	go test_do(fn, 200)
	go test_do(fn, 300)
	go sche.Handler()
	time.Sleep(time.Second * 100)
	fmt.Println("--------")
}

// 一个固定流程
func test_do(fn func(tasks []Task, final Task), start int) {
	fmt.Println("begin routine:", common.GetRoutineID())
	fn([]Task{func(args ...interface{}) {
		callback, _ := args[0].(Callback)
		//fmt.Println(callback)
		fmt.Println("task0 routine:", common.GetRoutineID())
		callback(false, start, 2)
	}, func(args ...interface{}) {
		callback, _ := args[0].(Callback)
		x, _ := args[1].(int)
		y, _ := args[2].(int)
		fmt.Println("task1 routine:", common.GetRoutineID())
		callback(false, x+y)
	}, func(args ...interface{}) {
		fmt.Println("task2 routine:", common.GetRoutineID())
		go func() {
			fmt.Println("enter newgo:", common.GetRoutineID())
			time.Sleep(time.Second * 2)
			callback, _ := args[0].(Callback)
			x, _ := args[1].(int)
			callback(false, x)
			fmt.Println("newgo over:", common.GetRoutineID())
		}()

	}}, func(args ...interface{}) {
		fmt.Println("final routine:", common.GetRoutineID())
		fmt.Println(args...)
	})
}
