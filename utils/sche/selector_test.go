package sche

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Selector(t *testing.T) {
	fmt.Println("---- Test_Selector ----")

	selectors := NewMultiSelector()

	go selectors.HandleLoop()

	time.Sleep(1 * time.Second)

	var chs1 = make(chan int)
	var chs2 = make(chan float64)
	var chs3 = make(chan string)

	var num1, num2, num3 int

	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs1),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}
			fmt.Println(v.Int())
			num1++
		}))

	// selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs1),
	// 	func(v reflect.Value, recvOk bool) {
	// 		if !recvOk {
	// 			return
	// 		}
	// 		fmt.Printf("do2 %v\n", v.Int())
	// 	}))

	// selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs1),
	// 	func(v reflect.Value, recvOk bool) {
	// 		if !recvOk {
	// 			return
	// 		}
	// 		fmt.Printf("do3 %v\n", v.Int())
	// 	}))

	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs2),
		func(v reflect.Value, recvOk bool) {
			fmt.Println(v.Float())
			num2++
		}))

	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs3),
		func(v reflect.Value, recvOk bool) {
			fmt.Println(v.String())
			num3++
		}))

	go func() {
		time.Sleep(3 * time.Second)
		//chs1 <- 9
	}()

	go func() {
		time.Sleep(1 * time.Second)
		fmt.Printf("close chs1\n")
		close(chs1)
	}()

	go func() {
		chs1 <- 9
		chs1 <- 8
		chs1 <- 7
		time.Sleep(2 * time.Second)
		chs2 <- 999
	}()

	go func() {
		time.Sleep(3 * time.Second)
		chs3 <- "hello"

	}()

	go func() {
		time.Sleep(3 * time.Second)
		chs2 <- 2
	}()

	go func() {
		time.Sleep(3 * time.Second)
		chs2 <- 3
	}()

	time.Sleep(5 * time.Second)

	// 有几个消息顺序随机

	// 确认参数
	//fmt.Printf("%v", results)
	assert.Equal(t, 3, num1, "missmatch number")
	assert.Equal(t, 3, num2, "missmatch number")
	assert.Equal(t, 1, num3, "missmatch number")
}
