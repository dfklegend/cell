
package common

type ReqIdType uint32

// zero代表无效
func IsValidReqId(id ReqIdType) bool {
    return id > 0
}