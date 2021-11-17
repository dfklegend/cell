module github.com/dfklegend/cell/server/examples/chat-client-go

go 1.16

replace github.com/dfklegend/cell/utils => ../../../utils

replace github.com/dfklegend/cell/net => ../../../net

replace github.com/dfklegend/cell/client => ../../../client

replace github.com/dfklegend/cell/server => ../../../server

replace github.com/dfklegend/cell/rpc => ../../../rpc

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/dfklegend/cell/client v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/net v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/utils v0.0.0
)
