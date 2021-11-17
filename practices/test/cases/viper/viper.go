package main


import (
	"fmt"
    "reflect"
	
    "github.com/spf13/viper"
)

type Server struct {
    ID string
    Type string
    Frontend bool
    ClientAddress string
    Address string
}

type Setting struct {
    Servers map[string]*Server
}

type Setting1 struct {
    Servers []*Server
}

func main() {
    //loadServer1()
    //loadServer2()
    //unmarshalKey()
    unmarshalKey1()
}

func loadServer1() {
    v := viper.New()
    v.SetConfigName("servers1")
    v.AddConfigPath("./config/")
    v.SetConfigType("yaml")
    v.ReadInConfig()    

    var obj Setting1

    v.Unmarshal(&obj)
    for i, one := range obj.Servers {
        fmt.Printf("%v %v\n", i, one)
        fmt.Printf("%v\n", reflect.TypeOf(one))
    }    

    fmt.Printf("%v\n", obj)
}

func loadServer2() {
	v := viper.New()
    v.SetConfigName("servers2")
    v.AddConfigPath("./config/")
    v.SetConfigType("yaml")
    v.ReadInConfig()    

    var obj Setting

    v.Unmarshal(&obj)
    for i, one := range obj.Servers {
        one.ID = i
        fmt.Printf("%v %v\n", i, one)
        fmt.Printf("%v\n", reflect.TypeOf(one))
    }    

    fmt.Printf("%v\n", obj)
}

func unmarshalKey() {
    v := viper.New()
    v.SetConfigName("servers")
    v.AddConfigPath("./config/")
    v.SetConfigType("yaml")
    v.ReadInConfig()    

    var server Server
    server.ID = "abcd"
    server.Type = "db"
    server.Frontend = false

    // 会覆盖掉定义数据
    v.UnmarshalKey("logic-1", &server)
    fmt.Printf("%v\n", server)
}

func unmarshalKey1() {
    v := viper.New()
    v.SetConfigName("servers")
    v.AddConfigPath("./config/")
    v.SetConfigType("yaml")
    v.ReadInConfig()    

    var server Server
    server.ID = "abcd"
    server.Type = "db"
    server.Frontend = false

    // 找不到key,不会出错
    v.UnmarshalKey("logic-2", &server)
    fmt.Printf("%v\n", server)
}
