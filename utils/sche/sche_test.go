package sche

import (
	"fmt"
	"testing"
	"time"

	"github.com/dfklegend/cell/utils/common"

	"github.com/stretchr/testify/assert"
)

type Null struct{}

// 测试处理中产生了panic
func Test_Panic(t *testing.T) {
	fmt.Println("---- Test_Panic ----")

	sche := NewSche()

	var rtSche, rtCall int
	waitOver := false
	go func() {
		fmt.Printf("%v call routine:%d\n", time.Now(), common.GetRoutineID())
		rtCall = common.GetRoutineID()
		t := sche.Post(func() (interface{}, error) {
			fmt.Printf("exec in routine:%d\n", common.GetRoutineID())
			rtSche = common.GetRoutineID()
			//var x interface{} = 7
			//return x.(struct{}), nil
			//var x = Null{}
			time.Sleep(3 * time.Second)
			panic("xxxx")
			//return 7, nil
		})
		// never got
		r, e := t.WaitDone()
		waitOver = true
		if e != nil {
			fmt.Println(e)
		}
		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
		fmt.Println(r)
	}()

	go sche.Handler()

	time.Sleep(8 * time.Second)

	// 检查结果
	assert.NotEqual(t, rtSche, rtCall, "应该不在同一个routine")
	assert.Equal(t, true, waitOver, "不会打断wait")
}

// 模拟目前一个死锁
func Test_Deadlock(t *testing.T) {
	fmt.Println("---- Test_Deadlock ----")

	sche := NewSche()

	waitt1Over := false
	go func() {
		fmt.Printf("%v call routine:%d\n", time.Now(), common.GetRoutineID())
		t := sche.Post(func() (interface{}, error) {
			fmt.Printf("exec in routine:%d\n", common.GetRoutineID())
			time.Sleep(3 * time.Second)

			fmt.Printf("post t1:%d\n", common.GetRoutineID())
			t1 := sche.Post(func() (interface{}, error) {
				fmt.Printf("t1:%d\n", common.GetRoutineID())
				return nil, nil
			})
			// will never end
			t1.WaitDone()
			fmt.Printf("%v", t1)
			waitt1Over = true
			return 7, nil
		})
		r, e := t.WaitDone()
		//fmt.Printf("%v", e)
		if e != nil {
			fmt.Println(e)
		}
		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
		fmt.Println(r)
	}()

	go sche.Handler()

	time.Sleep(8 * time.Second)
	assert.Equal(t, false, waitt1Over, "死锁")
}

func Test_ScheClose(t *testing.T) {
	fmt.Println("---- Test_ScheClose ----")

	sche := NewSche()
	fmt.Printf("%v main routine:%d\n", time.Now(), common.GetRoutineID())
	_ = sche.Post(func() (interface{}, error) {
		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
		return nil, nil
	})
	go sche.Handler()

	time.Sleep(1 * time.Second)

	sche.Stop()
	_ = sche.Post(func() (interface{}, error) {
		fmt.Printf("%v call routine over:%d\n", time.Now(), common.GetRoutineID())
		return nil, nil
	})

	time.Sleep(3 * time.Second)
}
