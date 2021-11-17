package config

import (
    "fmt"
    "testing"
)


func TestLoadConfig(t *testing.T) {
    config := LoadConfig("../data/config/")
    fmt.Printf("%v,%v,%v\n", config.GetInt("exam.int"),
        config.GetString("exam.int"), config.GetString("exam.string"))    

    fmt.Printf("%v\n", config.GetInt("exam.int2", 99))
    fmt.Printf("%v\n", config.GetInt("exam.int2"))

    fmt.Printf("%v\n", config.GetString("exam.int2", "def"))
    fmt.Printf("%v\n", config.GetDuration("exam.int2"))
}

