package apientry

import (
	"reflect"
)

// ----------------

// 代表调用session
type IHandlerSession interface {
	Reserve()
	Handle()
}

type DummySession struct {
}

func (self *DummySession) Reserve() {}
func (self *DummySession) Handle() {}

// API格式器
// 提供分析API，API调用抽象
// 前端接口和rpc
// receiver, [session], inArg, cb
// rpc接口没有session
// handler接口有session
type IAPIFormater interface {
	// 分析是否符合接口需求
	IsValidMethod(reflect.Method) bool
	NeedSession() bool
	// Receiver, finalArgs, cbFunc
	MakeCallArgs(reflect.Value, reflect.Value, reflect.Value, interface{}) []reflect.Value
}

var (
	TypeOfSession = reflect.TypeOf((*IHandlerSession)(nil)).Elem()
)

// API入口
type IAPIEntry interface {
	Desc() string
}

// ----------------
type BaseAPIEntry struct {
}

func (self *BaseAPIEntry) Desc() string {
	return "BaseAPIEntry"
}

type APIEntry struct {
	BaseAPIEntry
}
