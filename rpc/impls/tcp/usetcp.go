package tcpimpl

import (
    "github.com/dfklegend/cell/rpc/config"    
)

func Visit() {    
}

func UseTcp() {
    config.SetCreateClientImplFunc(NewClientImpl)
    config.SetCreateServerImplFunc(NewServerImpl)
}

