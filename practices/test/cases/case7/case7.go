package main


import (
	"fmt"
	//"reflect"
)

// ---------- 回调函数
type CBFunc func(e error, result interface{})


func TestFunc1(a int, cbFunc CBFunc) {
    cbFunc(nil, a + 100)
}

func TestFunc2(str string, cbFunc CBFunc) {
    TestFunc1(100, func(e error, result interface{}) {
        intV := result.(int)
        cbFunc(nil, fmt.Sprintf("%v%v", str, intV))
    })
}

func TestFunc() {
    TestFunc2("haha", func(e error, result interface{}) {
        str := result.(string)
        fmt.Println(str)
    })
}

// ----------- 回调对象
type ICBObj interface {
    Result(error, interface{})
}

type FuncCBObj struct {
    ResultCB  CBFunc
}

func (self *FuncCBObj) Result(e error, result interface{}) {
    self.ResultCB(e, result)
}

func NewFuncCBObj(cb CBFunc) ICBObj {
    return &FuncCBObj {
        ResultCB: cb,
    }
}

func TestObj1(a int, cbObj ICBObj) {
    cbObj.Result(nil, a + 100)
}

func TestObj2(str string, cbObj ICBObj) {
    TestObj1(100, NewFuncCBObj(func(e error, result interface{}) {
        intV := result.(int)
        //cbFunc(nil, fmt.Sprintf("%v%v", str, intV))
        cbObj.Result(nil, fmt.Sprintf("%v%v", str, intV))
    }) )
}

func TestObj() {
    v1 := "hello"
    TestObj2("haha", NewFuncCBObj(func(e error, result interface{}) {
        str := result.(string)
        v2 := fmt.Sprintf("%v%v", str, v1)
        fmt.Println(v2)
    }))
}


// 测试json
// interface{} json化


// 测试cb函数透传
func main() {
	//TestFunc()
    TestObj()
}
