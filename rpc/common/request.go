package common

import (
	"context"
	"reflect"
	"time"
	
	"github.com/dfklegend/cell/rpc/consts"
	"github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/logger"
	"github.com/dfklegend/cell/utils/sche"
)

var (
	typeOfError = consts.TypeOfError
)

var ReqIdService *common.SerialIdService = common.NewSerialIdService()

// 代表一个请求
type Request struct {
	ReqId ReqIdType

	Ctx context.Context

	Method string
	InArg  interface{}

	// 回调函数
	// 参数(error, arg)
	CBFunc reflect.Value
	// 通过CBFunc分析出来的
	// 用于回传消息时反序列化
	CBArgType reflect.Type
	// options
	Scheduler *sche.Sche

	// ack
	TimeOut int64
}

func NewRequest(method string,
	inArg interface{}, options ...ReqOption) *Request {
	req := &Request{
		ReqId:   ReqIdType(ReqIdService.AllocId()),
		Method:  method,
		InArg:   inArg,
		TimeOut: time.Now().Unix() + consts.RPCTimeout,
	}

	// apply options
	for _, opt := range options {
		opt(req)
	}

	req.processCB()
	return req
}

// 处理回调
// https://stackoverflow.com/questions/26321115/using-reflection-to-call-a-function-with-a-nil-parameter-results-in-a-call-usin
func (self *Request) processCB() {
	if !self.CBFunc.IsValid() {
		return
	}
	cbFunc := self.CBFunc
	cb := cbFunc.Interface()
	typ := reflect.TypeOf(cb)
	if typ.Kind() != reflect.Func {
		return
	}

	// outArg, error
	if typ.NumIn() != 2 {
		return
	}

	self.CBArgType = typ.In(0)
	// 检查第2个对象是否是error
	if typ.In(1) != typeOfError {
		logger.Log.Infof("cb second arg must be error")
	}
}

func (self *Request) NeedAck() bool {
	return self.CBFunc.IsValid()
}

// 动态配置
// -------------------------
type ReqOption func(*Request)

func ReqWithCB(cb reflect.Value) ReqOption {
	return func(self *Request) {
		self.CBFunc = cb
	}
}

func ReqWithScheName(name string) ReqOption {
	return func(self *Request) {
		self.Scheduler = sche.DefaultScheMgr.GetSche(name)
	}
}

func ReqWithSche(scheduler *sche.Sche) ReqOption {
	return func(self *Request) {
		self.Scheduler = scheduler
	}
}

// -------------------------
