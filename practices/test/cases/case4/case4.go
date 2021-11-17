package main

// 测试函数指针
 
import (	
	"fmt"
	"reflect"
)
 
type CBFunc1 func(interface{})
type CBFunc2 func(args ...interface{})

type Exam struct {
	ABC string;
}


func DumpFunc1(f CBFunc1, t reflect.Type) {
	fmt.Printf("%v", f)
}

func DumpFunc2(f CBFunc2, t reflect.Type) {
	fmt.Printf("%v", f)
}
 
// func main() {
// 	DumpFunc2(func(args ...interface{}) {
// 	}, reflect.TypeOf(Exam{}))
	
// 	var t reflect.Type	
// 	t = reflect.TypeOf(Exam{})
// 	fmt.Printf("%v", t)
// }
 

func main() {	
	// struct
	fmt.Printf("%v", reflect.TypeOf(Exam{}))
}


