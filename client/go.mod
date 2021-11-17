module github.com/dfklegend/cell/client

go 1.16

replace github.com/dfklegend/cell/utils => ../utils

replace github.com/dfklegend/cell/net => ../net

require (
	github.com/dfklegend/cell/net v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/utils v0.0.0
)
