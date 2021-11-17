package server

// RPC转接器
// 未来可以扩展不同
type RPCAdapter interface {
	Init(*RPCServerNode)
	Stop()
}
