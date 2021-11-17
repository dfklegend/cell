package service

import (
    "errors"
    "reflect"
    "fmt"
    "encoding/json"
)

// 对应一个函数
type Handler struct {
    // 对应的函数
    Method      reflect.Method
    // 参数1 receiver
    // 参数2类型
    ArgsType    reflect.Type;
}

type Service struct {
    Name      string              // name of service
    Type      reflect.Type        // type of the receiver
    Receiver  reflect.Value       // receiver of methods for the service
    // "someMethod": Handler
    Handlers  map[string]*Handler // registered methods
    Options   options             // options
}

func NewService(comp Component, opts []Option) *Service {
    s := &Service{
        Type:     reflect.TypeOf(comp),
        Receiver: reflect.ValueOf(comp),
    }

    // apply options
    for i := range opts {
        opt := opts[i]
        opt(&s.Options)
    }
    
    s.Name = reflect.Indirect(s.Receiver).Type().Name()
    return s
}

// suitableMethods returns suitable methods of typ
func (s *Service) suitableHandlerMethods(typ reflect.Type) map[string]*Handler {
    methods := make(map[string]*Handler)
    for m := 0; m < typ.NumMethod(); m++ {
        method := typ.Method(m)
        mt := method.Type
        mn := method.Name
        if isHandlerMethod(method) { 
            methods[mn] = &Handler{Method: method, ArgsType: mt.In(2)}
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
func (s *Service) ExtractHandler() error {
    typeName := reflect.Indirect(s.Receiver).Type().Name()
    if typeName == "" {
        return errors.New("no service name for type " + s.Type.String())
    }
    if !isExported(typeName) {
        return errors.New("type " + typeName + " is not exported")
    }

    // Install the methods
    s.Handlers = s.suitableHandlerMethods(s.Type)

    if len(s.Handlers) == 0 {
        str := ""
        // To help the user, see if a pointer receiver would work.
        method := s.suitableHandlerMethods(reflect.PtrTo(s.Type))
        if len(method) != 0 {
            str = "type " + s.Name + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
        } else {
            str = "type " + s.Name + " has no exported methods of suitable type"
        }
        return errors.New(str)
    }

    for i := range s.Handlers {
        one := s.Handlers[i]
        fmt.Printf("found Handler:%v %v\n", i, one)
    }

    return nil
}

// 查找
func (s *Service) CallMethod(method string, msg []byte) error {
    handler := s.Handlers[method]
    if handler == nil {
        return fmt.Errorf("can not find method:%v", method)
    }
    // call it
    // 转化参数
    finalArgs := reflect.New(handler.ArgsType.Elem()).Interface()
    err := json.Unmarshal(msg, finalArgs)
    if err != nil {
        return err
    }
    // 
    args := []reflect.Value{s.Receiver, reflect.ValueOf(1), reflect.ValueOf(finalArgs)}
    handler.Method.Func.Call(args)
    return nil
}