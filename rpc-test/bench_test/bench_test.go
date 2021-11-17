package client

import (	
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/rpc/common"
	"github.com/dfklegend/cell/rpc/client"
	"github.com/dfklegend/cell/rpc/server"
	"github.com/dfklegend/cell/rpc/test"	
)

var srv *server.RPCServerNode
var ms *client.MailStation

func Benchmark_Normal(b *testing.B) {	
	if srv == nil {
		srv = server.NewNode("node1")
		srv.Start(":50000")

		srv.GetCollection().Register(&test.Entry1{}, api.WithNameFunc(strings.ToLower))
		srv.GetCollection().Build()

		ms = client.NewMailStation("self")
		ms.AddServer("s0", fmt.Sprintf("localhost:%v", 50000))

		time.Sleep(3 * time.Second)
	}

	b.ResetTimer()

	log.Println("bench start")
	for i := 0; i < b.N; i++ {
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.join",
			&test.InMsg{fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				//log.Printf("rpc got result:%v %v\n", result, e)
			})))
	}
	log.Println("call over")
	waitMatch()
	log.Println("wait over")
	ms.DumpStat()

	//srv.Stop()
}

func waitMatch() {
	for !ms.IsAllBenchOver() {
		time.Sleep(1 * time.Second)
	}
}

// go test -benchmem -bench ^Benchmark_Normal -benchtime=10000x
// 测试吞吐量相关参数
// pullAcks的频率和一次结果次数
// 50ms 400
// 100000 146329 ns/op
// 近似5000次/s

// 调整为单routine后，
// 测试10w次，出现了sche 队列满的情况
// 增加队列上限 99999
// 100000 161302 ns/op 10323 B/op 217 allocs/op

// 使用TCP协议
// 100000 321359 ns/op 808057 B/op 104 allocs/op
// 
// 之前不小心开了压缩，关掉后
// 使用TCP协议
// 100000 75985 ns/op 4618 B/op 77 allocs/op