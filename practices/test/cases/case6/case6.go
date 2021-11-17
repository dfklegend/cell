package main

// 测试reflect

import (
	"fmt"
	"reflect"
)

func TestFunc() {
}

type Exam struct {
}

func (self *Exam) TestFunc1() {
}

func main() {
	t := reflect.TypeOf(TestFunc)
	fmt.Printf("%v %v %v\n", t, t.Kind(), reflect.Func)

	var v reflect.Value
	fmt.Printf("%v %v\n", v, v.IsValid())

	fmt.Printf("-----------------")
}
