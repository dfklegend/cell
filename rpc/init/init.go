package initp

import (
    "github.com/dfklegend/cell/rpc/impls/tcp"
    "github.com/dfklegend/cell/rpc/impls/grpc"
)

// 缺省使用TCP
var DefaultUseTCP = true

func init() {
    if DefaultUseTCP {
        tcpimpl.UseTcp()    
    } else {
        grpcimpl.UseGRPC()
    }    
}

func Visit() {    
}
