package common

import (
    "fmt"
    "net/http"
    _ "net/http/pprof"	
)

func GoPprofServe(port string) {
    go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", port), nil)
}

