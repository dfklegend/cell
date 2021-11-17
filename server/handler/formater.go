package handler

import (
	"reflect"

	//"github.com/dfklegend/cell/utils/logger"
	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/rpc/consts"
)

type HandlerFormater struct {
	dummy api.DummySession
}

func NewHandlerFormater() api.IAPIFormater {
	return &HandlerFormater{}
}

func (self *HandlerFormater) NeedSession() bool {
	return true
}

func (self *HandlerFormater) IsValidMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs three ins: receiver, *Session, []byte or pointer, cb.
	if mt.NumIn() != 4 {
		//log.Printf("%v must had %v args\n", method.Name, 4)
		return false
	}

	// Method needs one outs: error
	if mt.NumOut() != 1 {
		return false
	}

	// t1是*session	
	if t1 := mt.In(1); !t1.Implements(api.TypeOfSession) {		
		return false
	}	

	if (mt.In(2).Kind() != reflect.Ptr && mt.In(2) != consts.TypeOfBytes) || mt.Out(0) != consts.TypeOfError {
		return false
	}
	return true
}

func (self *HandlerFormater) MakeHandlerSession(from interface{}) api.IHandlerSession {
	return from.(api.IHandlerSession)
}

func (self *HandlerFormater) MakeCallArgs(receiver reflect.Value, args reflect.Value,
	cbFunc reflect.Value, ext interface{}) []reflect.Value {
	session := self.MakeHandlerSession(ext)
	return []reflect.Value{receiver, reflect.ValueOf(session), args, cbFunc}
}

