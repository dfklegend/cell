package grpcimpl

import (
    //"github.com/dfklegend/cell/rpc/server"    
    "github.com/dfklegend/cell/rpc/interfaces"
)

// grpc实现
type GRPCServerImpl struct {  
    node interfaces.IRPCNode
    server *GRPCServer
}

func NewServerImpl() interfaces.IRPCServerImpl {
    return &GRPCServerImpl{}
}

func (self *GRPCServerImpl) Init(n interfaces.IRPCNode) {    
    self.node = n
}

func (self *GRPCServerImpl) Start(listenAddress string) {
    server := NewGRPCServer()
    server.Init(self.node)
    server.Start(listenAddress)
    self.server = server
}

func (self *GRPCServerImpl) Stop() {
    self.server.Stop()
}

