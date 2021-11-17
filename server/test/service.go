package test

import (
	"errors"
	"fmt"
	"math/rand"

	api "github.com/dfklegend/cell/rpc/apientry"
)

// 接口
type HandlerEntry struct {
	api.APIEntry
}

func (self *HandlerEntry) Join(d *api.DummySession, msg *InMsg, cbFunc api.HandlerCBFunc) error {

	str := fmt.Sprintf("%v->%v", msg.Abc, rand.Intn(100))
	//log.Printf("in join:%v\n", str)

	api.CheckInvokeCBFunc(cbFunc, nil,
		&OutMsg{str})
	return nil
}

// 错误返回流程
func (self *HandlerEntry) Error(d *api.DummySession, msg *InMsg, cbFunc api.HandlerCBFunc) error {
	api.CheckInvokeCBFunc(cbFunc, errors.New("some error"),
		&OutMsg{})
	return nil
}

func (self *HandlerEntry) Panic(d *api.DummySession, msg *InMsg, cbFunc api.HandlerCBFunc) error {
	panic("hahaha")
	api.CheckInvokeCBFunc(cbFunc, errors.New("some error"),
		&OutMsg{})
	return nil
}
