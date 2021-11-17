module github.com/dfklegend/cell/rpc

go 1.16

require (
	github.com/stretchr/testify v1.7.0
	google.golang.org/grpc v1.40.0
)

require (
	github.com/dfklegend/cell/client v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/net v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/utils v0.0.0
	github.com/golang/protobuf v1.5.2
)

replace github.com/dfklegend/cell/utils => ../utils

replace github.com/dfklegend/cell/client => ../client

replace github.com/dfklegend/cell/net => ../net
