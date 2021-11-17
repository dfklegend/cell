package apientry

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"log"

	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/logger"
	"github.com/dfklegend/cell/utils/sche"
)

type HandlerCBFunc func(error, interface{})

// 对应一个函数
type Handler struct {
	// 对应的函数
	// 函数参数 receiver, arg, handlerCBFunc
	Method reflect.Method
	// 参数1 receiver
	// 参数2类型
	ArgsType reflect.Type
}

// 基于入口对象分析构建而成
type APIContainer struct {
	Name     string        // name of service
	Type     reflect.Type  // type of the receiver
	Receiver reflect.Value // receiver of methods for the service
	// "someMethod": Handler
	Handlers map[string]*Handler // registered methods
	Options  options             // options

	schedName string // 调度器名字
}

// TODO: 目标函数的格式和调用 抽象出来
// 便于 远程接口和客户端接口使用
func NewContainer(entry IAPIEntry, opts ...Option) *APIContainer {
	s := &APIContainer{
		Type:     reflect.TypeOf(entry),
		Receiver: reflect.ValueOf(entry),
	}

	// apply options
	for _, opt := range opts {
		opt(&s.Options)
	}

	if name := s.Options.name; name != "" {
		s.Name = name
	} else {
		s.Name = reflect.Indirect(s.Receiver).Type().Name()
		if s.Options.nameFunc != nil {
			s.Name = s.Options.nameFunc(s.Name)
		}
	}
	s.schedName = s.Options.schedName

	return s
}

// suitableMethods returns suitable methods of typ
func (self *APIContainer) suitableHandlerMethods(formater IAPIFormater, typ reflect.Type) map[string]*Handler {
	methods := make(map[string]*Handler)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mt := method.Type
		mn := method.Name

		// rewrite handler name
		if self.Options.nameFunc != nil {
			mn = self.Options.nameFunc(mn)
		}

		if formater != nil && formater.IsValidMethod(method) {
			methods[mn] = &Handler{Method: method, ArgsType: mt.In(FormaterGetInArgIndex(formater))}
		}
	}
	return methods
}

// ExtractHandler extract the set of methods from the
// receiver value which satisfy the following conditions:
// - exported method of exported type
// - two arguments, both of exported type
// - the first argument is *session.Session
// - the second argument is []byte or a pointer
func (self *APIContainer) ExtractHandler(formater IAPIFormater) error {
	typeName := reflect.Indirect(self.Receiver).Type().Name()
	if typeName == "" {
		return errors.New("no service name for type " + self.Type.String())
	}
	if !isExported(typeName) {
		return errors.New("type " + typeName + " is not exported")
	}

	// Install the methods
	self.Handlers = self.suitableHandlerMethods(formater, self.Type)

	if len(self.Handlers) == 0 {
		str := ""
		// To help the user, see if a pointer receiver would work.
		method := self.suitableHandlerMethods(formater, reflect.PtrTo(self.Type))
		if len(method) != 0 {
			str = "type " + self.Name + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			str = "type " + self.Name + " has no exported methods of suitable type"
		}
		return errors.New(str)
	}

	for i := range self.Handlers {
		one := self.Handlers[i]
		logger.Log.Infof("found Handler:%v %v\n", i, one)
	}

	return nil
}

func (self *APIContainer) HasMethod(method string) bool {
	return self.Handlers[method] != nil
}

// 查找Handler
// 调用具体函数，会将[]byte数据反序列化成参数结构
// 返回的数据结构也会被序列化成[]byte
/**
 * @param formater{IAPIFormater} 用于提供接口类型差异化比如handler和remote
 * @param method{string} 方法名
 * @param msg{[]byte} 此接口负责反序列化成接口参数结构
 * @param cbFunc{HandlerCBFunc}
 *        cb(e error, outArg interface{})
 *        	e 错误
 *        	outArg 接口返回的对象
 *        	(此函数负责序列化，返回出去)
 * @param ext{interface{}} 见collection说明
 */
func (self *APIContainer) CallMethod(formater IAPIFormater, method string,
	msg []byte, cbFunc HandlerCBFunc, ext interface{}) error {
	//log.Printf("enter: %v", method)
	handler := self.Handlers[method]
	if handler == nil {
		return fmt.Errorf("can not find method:%v", method)
	}
	// call it
	// 转化为目标参数
	finalArgs := reflect.New(handler.ArgsType.Elem()).Interface()
	err := json.Unmarshal(msg, finalArgs)
	if err != nil {
		fmt.Printf("arg json.Unmarshal [%v] failed: %v\n", string(msg), err)
		//return err
	}

	replaceCB := func(e error, result interface{}) {
		// 序列化返回值		
		data, jerr := json.Marshal(result)
		e1 := jerr
		if e1 == nil {
			e1 = e
		}
		CheckInvokeCBFunc(cbFunc, e1, data)
	}

	// 参数列表
	// args := []reflect.Value{self.Receiver, reflect.ValueOf(1),
	// 	reflect.ValueOf(finalArgs), reflect.ValueOf(replaceCB)}
	args := formater.MakeCallArgs(self.Receiver, reflect.ValueOf(finalArgs), reflect.ValueOf(replaceCB), ext)

	//log.Printf("call: %v", method)
	// 根据schedName 调度到指定的
	if self.schedName != "" {
		scheduler := sche.DefaultScheMgr.GetSche(self.schedName)
		if scheduler == nil {
			logger.Log.Errorf("can not find sche:%v", self.schedName) 
			return nil
		}
		scheduler.Post(func() (interface{}, error) {
			SafeCall(handler, args, cbFunc)
			return nil, nil
		})
		return nil
	}
	SafeCall(handler, args, cbFunc)
	return nil
}

// 捕获异常
func SafeCall(handler *Handler, args []reflect.Value, cbFunc HandlerCBFunc) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Infof("panic in handler.call:%v", err)

			stack := common.GetStackStr()
			log.Printf(stack)
			logger.Log.Infof(stack)
			CheckInvokeCBFunc(cbFunc, errors.New("panic in rpc"), nil)
		}
	}()

	handler.Method.Func.Call(args)
}

func CheckInvokeCBFunc(cbFunc HandlerCBFunc, e error, result interface{}) {
	if cbFunc == nil {
		return
	}
	cbFunc(e, result)
}
