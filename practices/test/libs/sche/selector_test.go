package sche

import (
	"fmt"
	"testing"
	"time"
	"reflect"
)

func Test_Selector(t *testing.T) {
	fmt.Println("---- Test_Selector ----")

	selectors := NewMultiSelector();

	
	go selectors.HandleLoop();

	time.Sleep(1 * time.Second)

	var chs1 = make(chan int)
    var chs2 = make(chan float64)
	var chs3 = make(chan string)
	
	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs1), 
		func(v reflect.Value, recvOk bool) {
			if(!recvOk) {
				return
			}			
			fmt.Println(v.Int())
		}) )

	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs1), 
		func(v reflect.Value, recvOk bool) {
			if(!recvOk) {
				return
			}			
			fmt.Printf("do2 %v\n", v.Int())
		}) )

	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs2), 
		func(v reflect.Value, recvOk bool) {
			fmt.Println(v.Float())
		}) )

	selectors.AddSelector(NewFuncSelector(reflect.ValueOf(chs3), 
		func(v reflect.Value, recvOk bool) {
			fmt.Println(v.String())
		}) )	

	go func() {
		time.Sleep(3 * time.Second)
		//chs1 <- 9
	}()

	go func() {
		time.Sleep(1 * time.Second)		
		close(chs1)
	}()

	go func() {
		chs1 <- 9
		time.Sleep(2 * time.Second)
		chs2 <- 999		
	}()

	go func() {
		time.Sleep(2 * time.Second)
		chs3 <- "hello"
		chs2 <- 2		
		chs2 <- 3		
	}()

	time.Sleep(5 * time.Second)

	// 确认参数
}


