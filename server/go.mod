module github.com/dfklegend/cell/server

go 1.16

require go.etcd.io/etcd v0.0.0-20210226220824-aa7126864d82

replace github.com/dfklegend/cell/utils => ../utils

require (
	github.com/dfklegend/cell/net v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/rpc v0.0.0-00010101000000-000000000000
	github.com/dfklegend/cell/utils v0.0.0
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/google/uuid v1.1.4
	github.com/kr/text v0.2.0 // indirect
	github.com/nats-io/nats-server/v2 v2.1.2
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
)

replace github.com/dfklegend/cell/rpc => ../rpc

replace github.com/dfklegend/cell/net => ../net

//replace google.golang.org/grpc/naming => ../../packages/grpc@v1.26.0\naming
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
