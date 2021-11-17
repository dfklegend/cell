package config

import (
    "github.com/dfklegend/cell/rpc/interfaces"
)


// 创建器 提供RPC底层实现的定制
type CreateClientImplFunc func() interfaces.IRPCClientImpl
type CreateServerImplFunc func() interfaces.IRPCServerImpl

var curCreateClientImplFunc CreateClientImplFunc
var curCreateServerImplFunc CreateServerImplFunc

func SetCreateClientImplFunc( f CreateClientImplFunc) {
    curCreateClientImplFunc = f
}

func CreateClientImpl() interfaces.IRPCClientImpl {
    return curCreateClientImplFunc()
}

func SetCreateServerImplFunc( f CreateServerImplFunc) {
    curCreateServerImplFunc = f
}

func CreateServerImpl() interfaces.IRPCServerImpl {
    return curCreateServerImplFunc()
}
