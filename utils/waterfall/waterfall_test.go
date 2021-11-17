package waterfall

import (
	"fmt"
	"testing"
	"time"

	"github.com/dfklegend/cell/utils/common"
)

func FTest_Simple(t *testing.T) {
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
	fmt.Println("begin routine:", common.GetRoutineID())
	fn([]Task{func(args ...interface{}) {
		callback, _ := args[0].(Callback)
		//fmt.Println(callback)
		fmt.Println("task0 routine:", common.GetRoutineID())
		callback(true, 1, 2)
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
