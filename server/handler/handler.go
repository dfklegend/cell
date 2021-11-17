package handler

import (
	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/utils/logger"
)

type HandlerNode struct {
	// 接口集
	collection *api.APICollection
}

func NewHandlerNode() *HandlerNode {
	h := &HandlerNode{
		collection: api.NewCollection(),
	}
	h.Init()
	return h
}

func (self *HandlerNode) GetCollection() *api.APICollection {
	return self.collection
}

// 设置合适的formater
func (self *HandlerNode) Init() {
	self.collection.SetFormater(NewHandlerFormater())
}

func (self *HandlerNode) Start() {
	self.collection.Build()
}

//
func (self *HandlerNode) Call(session api.IHandlerSession, method string, msg []byte, cbFunc api.HandlerCBFunc) {
	
	logger.Log.Debugf("Call:%v %v", method, string(msg))
	callerr := self.collection.Call(method, msg, func(e error, result interface{}) {
		cbFunc(e, result)
	}, session)

	if callerr != nil {
		cbFunc(callerr, nil)
	}
}
