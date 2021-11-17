module github.com/dfklegend/cell/rpc-test

go 1.16

replace github.com/dfklegend/cell/utils => ../utils

replace github.com/dfklegend/cell/rpc => ../rpc

replace github.com/dfklegend/cell/client => ../client

replace github.com/dfklegend/cell/net => ../net

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/dfklegend/cell/rpc v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/utils v0.0.0
	github.com/stretchr/testify v1.7.0
)
