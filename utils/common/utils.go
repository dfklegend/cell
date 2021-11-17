package common

import (
    "time"
	"runtime"
    "math/rand"
    "encoding/json"	
)

func GetStackStr() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
}

func SafeJsonMarshalByteArray(data interface{}) []byte {
    r, e := json.Marshal(data)
    if e != nil {
        return []byte("")
    }
    return r
}

func SafeJsonMarshal(data interface{}) string {
    r, e := json.Marshal(data)
    if e != nil {
        return ""
    }
    return string(r)
}

func NowMs() int64 {
    return time.Now().UnixNano()/1e6
}

func RandFloat32(min, max float32) float32 {
    return min + rand.Float32()*(max-min)
}