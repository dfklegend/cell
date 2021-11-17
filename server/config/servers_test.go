package config

import (
    "fmt"
    "testing"
)


func TestLoadServers(t *testing.T) {
    servers := LoadServers("../data/config/")
    fmt.Printf("%v\n", servers)
    for k, v := range servers.Servers {
        fmt.Printf("%v %v\n", k, v)
    }
}

