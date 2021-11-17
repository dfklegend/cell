package consts

import "errors"

// Errors that can occur during message handling.
var (
	ErrExample = errors.New("error example")
	// 请求处理不过来，堆积太多
	ErrRPCServerTooBusy = errors.New("error ErrRPCServerTooBusy")
	ErrRPCAutoTimeout   = errors.New("error rpc auto timeout")
)
