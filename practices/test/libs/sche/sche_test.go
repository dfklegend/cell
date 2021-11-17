package sche

import (
	"fmt"
	"testing"
	"time"

	"github.com/dfklegend/cell/practices/test/libs/util"
)

type Null struct{}

// . 建立一个sche
// . 建立一些子routine
func Test_Simple(t *testing.T) {
	fmt.Println("---- Test_Simple ----")

	sche := NewSche()

	go func() {
		fmt.Printf("%v call routine:%d\n", time.Now(), util.GetRoutineID())
		t := sche.Post(func() (interface{}, error) {
			fmt.Printf("exec in routine:%d\n", util.GetRoutineID())
			//var x interface{} = 7
			//return x.(struct{}), nil
			//var x = Null{}
			time.Sleep(3 * time.Second)
			panic("xxxx")
			//return 7, nil
		})
		r, e := t.WaitDone()
		//fmt.Printf("%v", e)
		if e != nil {
			fmt.Println(e)
		}
		fmt.Printf("%v call routine over:%d\n", time.Now(), util.GetRoutineID())
		fmt.Println(r)
	}()

	go sche.Handler()

	time.Sleep(8 * time.Second)
}

func Test_Simple1(t *testing.T) {
	fmt.Println("---- Test_Simple ----")

	sche := NewSche()

	go func() {
		fmt.Printf("%v call routine:%d\n", time.Now(), util.GetRoutineID())
		t := sche.Post(func() (interface{}, error) {
			fmt.Printf("exec in routine:%d\n", util.GetRoutineID())
			time.Sleep(3 * time.Second)

			fmt.Printf("post t1:%d\n", util.GetRoutineID())
			sche.Post(func() (interface{}, error) {
				fmt.Printf("t1:%d\n", util.GetRoutineID())
				return nil, nil
			})
			//t1.WaitDone()
			return 7, nil
		})
		r, e := t.WaitDone()
		//fmt.Printf("%v", e)
		if e != nil {
			fmt.Println(e)
		}
		fmt.Printf("%v call routine over:%d\n", time.Now(), util.GetRoutineID())
		fmt.Println(r)
	}()

	go sche.Handler()

	time.Sleep(8 * time.Second)
}
