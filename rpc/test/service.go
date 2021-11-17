package test

import (
	"errors"
	"fmt"
	"math/rand"

	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/utils/logger"
)

type Entry1 struct {
	api.APIEntry
}

func (self *Entry1) Join(msg *InMsg, cbFunc api.HandlerCBFunc) error {

	str := fmt.Sprintf("%v->%v", msg.Abc, rand.Intn(100))
	//log.Printf("in join:%v\n", str)

	api.CheckInvokeCBFunc(cbFunc, nil,
		&OutMsg{str})
	return nil
}

func (self *Entry1) Echo(msg *InMsg, cbFunc api.HandlerCBFunc) error {
	logger.Log.Debugf("server got echo:%v\n", msg.Abc)
	api.CheckInvokeCBFunc(cbFunc, nil,
		&OutMsg{msg.Abc})
	return nil
}

// 错误返回流程
func (self *Entry1) Error(msg *InMsg, cbFunc api.HandlerCBFunc) error {
	api.CheckInvokeCBFunc(cbFunc, errors.New("some error"),
		&OutMsg{})
	return nil
}

func (self *Entry1) Panic(msg *InMsg, cbFunc api.HandlerCBFunc) error {
	panic("hahaha")
	api.CheckInvokeCBFunc(cbFunc, errors.New("some error"),
		&OutMsg{})
	return nil
}
