package waterfall

import (
	"fmt"
	"log"

	"github.com/dfklegend/cell/utils/common"
)

type Callback func(err bool, args ...interface{})
type Task func(args ...interface{})

// 实现一个类似nodejs waterfall
// 每个回调都在调用routine中调用

// 简版
// 测试函数指针等
// just callback directly
func Waterfall_Simple(tasks []Task, final Task) {
	// 执行索引
	cursor := 0
	size := len(tasks)

	var callback Callback = nil
	var gofinal Task = nil

	exec := func(index int, args ...interface{}) {
		args = append([]interface{}{callback}, args...)
		tasks[index](args...)
	}

	gonext := func(args ...interface{}) {
		cursor++
		if cursor < size {
			exec(cursor, args...)
		} else {
			gofinal(args...)
		}
	}

	gofinal = func(args ...interface{}) {
		//args = append([]interface{}{false}, args...)
		final(args...)
	}

	callback = func(err bool, args ...interface{}) {
		//fmt.Printf("callback:%v\n", util.GetRoutineID())
		if err {
			// go final
			args = append([]interface{}{false}, args...)
			gofinal(args...)
			return
		}

		// go next
		gonext(args...)
	}

	exec(0)
}

// 每个Task都在调用者routine中执行
// every task will execute in one routine(caller routine)
func Waterfall_Go(tasks []Task, final Task) {
	// 执行索引
	cursor := 0
	size := len(tasks)
	//
	task_next := 0
	task_final := 1

	chanNext := make(chan int, 1)

	var callback Callback = nil
	var gofinal Task = nil
	var curArgs []interface{}

	exec := func(index int, args ...interface{}) {
		args = append([]interface{}{callback}, args...)
		tasks[index](args...)
	}

	gonext := func(args ...interface{}) {
		curArgs = args
		chanNext <- task_next
	}

	donext := func(args ...interface{}) {
		cursor++
		if cursor < size {
			exec(cursor, args...)
		} else {
			gofinal(args...)
		}
	}

	gofinal = func(args ...interface{}) {
		curArgs = args
		chanNext <- task_final
	}

	dofinal := func(args ...interface{}) {		
		final(args...)
		close(chanNext)
	}

	// 不在主routine中
	//
	callback = func(err bool, args ...interface{}) {
		log.Printf("callback:%v  %v\n", common.GetRoutineID(), args)
		if err {
			// go final
			args = append([]interface{}{true}, args...)

			log.Printf("%v\n", args)
			gofinal(args...)
			return
		}

		// go next
		gonext(args...)
	}

	// 启动
	exec(0)

	// 问题，会阻塞调用者的routine
	// wait for next task
	for {
		// 0: next
		// 1: final
		data, ok := <-chanNext
		if !ok {
			// over
			fmt.Println("waterfall over")
			return
		}

		if data == task_next {
			donext(curArgs...)
		} else {
			dofinal(curArgs...)
		}
	}
}
