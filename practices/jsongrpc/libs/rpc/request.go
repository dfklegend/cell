
package rpc

type Request struct {
	// 请求id,到服务器时分配
	// 本节点唯一，用于对应客户端的请求
	ReqId int
	// 客户端请求Id
	// 返回时附带
	ClientReqId int
	
}