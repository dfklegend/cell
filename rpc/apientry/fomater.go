package apientry

import (
	"reflect"
)

var defaultFormater IAPIFormater = &DefaultFormater{}

// 缺省的rpc格式
type DefaultFormater struct {
	dummy DummySession
}

func (self *DefaultFormater) NeedSession() bool {
	return false
}

// handler函数 3个参数
// 		receiver, pointer, cb
func (self *DefaultFormater) IsValidMethod(method reflect.Method) bool {
	mt := method.Type
	// Method must be exported.
	if method.PkgPath != "" {
		return false
	}

	// Method needs three ins: receiver, pointer, cb.
	if mt.NumIn() != 3 {
		//log.Printf("%v must had %v args\n", method.Name, 3)
		return false
	}

	// Method needs one outs: error
	if mt.NumOut() != 1 {
		return false
	}

	// // t1是*session
	// if t1 := mt.In(1); t1.Kind() != reflect.Ptr || t1.Elem().Implements(TypeOfSession) {
	// 	return false
	// }

	// 参数1，返回
	if (mt.In(1).Kind() != reflect.Ptr && mt.In(1) != typeOfBytes) || mt.Out(0) != typeOfError {
		return false
	}
	return true
}

func (self *DefaultFormater) MakeCallArgs(receiver reflect.Value, args reflect.Value,
	cbFunc reflect.Value, ext interface{}) []reflect.Value {
	return []reflect.Value{receiver, args, cbFunc}
}

func FormaterGetInArgIndex(formater IAPIFormater) int {
	if formater.NeedSession() {
		return 2
	}
	return 1
}

