package waterfall

import (
	"fmt"
	"testing"
	"time"

	"github.com/dfklegend/cell/practices/test/libs/util"
)

func Test_Simple(t *testing.T) {
	fmt.Println("---- show waterfall callback ----")
	go test_do(Waterfall_Simple)
	time.Sleep(time.Second * 10)
	fmt.Println("--------")
}

func Test_Go(t *testing.T) {
	fmt.Println("---- show waterfall go ----")
	go test_do(Waterfall_Go)
	time.Sleep(time.Second * 10)
	fmt.Println("--------")
}

// 一个固定流程
func test_do(fn func(tasks []Task, final Task)) {
	fmt.Println("begin routine:", util.GetRoutineID())
	fn([]Task{func(args ...interface{}) {
		callback, _ := args[0].(Callback)
		//fmt.Println(callback)
		fmt.Println("task0 routine:", util.GetRoutineID())
		callback(false, 1, 2)
	}, func(args ...interface{}) {
		callback, _ := args[0].(Callback)
		x, _ := args[1].(int)
		y, _ := args[2].(int)
		fmt.Println("task1 routine:", util.GetRoutineID())
		callback(false, x+y)
	}, func(args ...interface{}) {
		fmt.Println("task2 routine:", util.GetRoutineID())
		go func() {
			fmt.Println("enter newgo:", util.GetRoutineID())
			time.Sleep(time.Second * 2)
			callback, _ := args[0].(Callback)
			x, _ := args[1].(int)
			callback(false, x)
			fmt.Println("newgo over:", util.GetRoutineID())
		}()

	}}, func(args ...interface{}) {
		fmt.Println("final routine:", util.GetRoutineID())
		fmt.Println(args...)
	})
}
