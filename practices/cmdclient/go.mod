module github.com/dfklegend/cell/cmdclient

go 1.16

require (
	github.com/dfklegend/cell/utils v0.0.0
	github.com/google/pprof v0.0.0-20210804190019-f964ff605595 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20210724235854-665d3a6fe486 // indirect
	github.com/sirupsen/logrus v1.0.5
	github.com/topfreegames/pitaya v0.0.0
	golang.org/x/sys v0.0.0-20210806184541-e5e7981a1069 // indirect
)

replace github.com/dfklegend/cell/utils => ../utils

replace github.com/topfreegames/pitaya => ../pitaya-notes-master
