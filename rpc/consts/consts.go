package consts

import (
	"reflect"
)

var (
	// Error的类型
	TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
	TypeOfBytes = reflect.TypeOf(([]byte)(nil))

	// 空error
	NilError = reflect.New(TypeOfError).Elem()
)

const (
	// RPC过期时间
	RPCTimeout = 30
)

const (
	DefaultServerPort = ":50051"
)
