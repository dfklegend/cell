package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"

	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/rpc/common"	
)

// node1
type Service1 struct {
	api.APIEntry
}

func (self *Service1) Func1(msg *InMsg, cbFunc api.HandlerCBFunc) error {
	api.CheckInvokeCBFunc(cbFunc, nil,
		&OutMsg{"service1.func1"})
	return nil
}

type Service2 struct {
	api.APIEntry
}

func (self *Service2) Func1(msg *InMsg, cbFunc api.HandlerCBFunc) error {

	// call node2 server3.func1
	ms1.Call("node2", "service3.func1",
		&InMsg{fmt.Sprintf("%v", rand.Intn(100))},
		common.ReqWithCB(reflect.ValueOf(func(result *OutMsg, e error) {
			log.Printf("Service2.func1 got result:%v %v\n", result, e)

			api.CheckInvokeCBFunc(cbFunc, nil,
				&OutMsg{result.Def})
		})))

	return nil
}

// node2
type Service3 struct {
	api.APIEntry
}

func (self *Service3) Func1(msg *InMsg, cbFunc api.HandlerCBFunc) error {
	api.CheckInvokeCBFunc(cbFunc, nil,
		&OutMsg{"service3.func1"})
	return nil
}
