package grpcimpl

import (
    "github.com/dfklegend/cell/rpc/config"    
)

func Visit() {    
}

func UseGRPC() {
    config.SetCreateClientImplFunc(NewClientImpl)
    config.SetCreateServerImplFunc(NewServerImpl)
}