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
	"github.com/dfklegend/cell/rpc/server"
	"github.com/dfklegend/cell/rpc/test"		
)

func Test_Init(t *testing.T) {	
}

// 正确和错误目标
func _Test_ServerRPC(t *testing.T) {
	

	log.Printf("--------- Test_ServerRPC ----------")
	do_ServerRPC(true)
}

// 测试RPC timeout处理
func _Test_RPCTimeout(t *testing.T) {
	log.Printf("--------- FuncTest_RPCTimeout ----------")
	flags := common.GetDebugFlags()
	flags.RPCNotSend = true
	do_ServerRPC(false)

	time.Sleep(45 * time.Second)
	flags.RPCNotSend = false
}

func do_ServerRPC(stopMs bool) {
	srv := server.NewNode("node1")
	srv.Start(":50000")

	srv.GetCollection().Register(&test.Entry1{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Build()

	time.Sleep(3 * time.Second)
	ms := NewMailStation("self")
	ms.AddServer("s0", fmt.Sprintf("localhost:%v", 50000))

	time.Sleep(3 * time.Second)

	sendRPC := true
	go func() {
		// 随机测试
		t := time.NewTicker(1 * time.Millisecond)
		
		defer t.Stop()		

		for sendRPC {
			select {
			case <-t.C:
				testCall(ms)
			}
		}
	}()

	time.Sleep(5 * time.Second)
	sendRPC = false

	time.Sleep(1 * time.Second)
	if stopMs {
		ms.Stop()	
	}
	
	srv.Stop()

	log.Printf("--------- stop ----------")
}

func testCall(ms *MailStation) {
	for i := 0; i < 10; i++ {
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.join",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				log.Printf("rpc got result:%v\n", result)
			})))
		//ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.join", "ddd")
	}
}

func _Test_Null(t *testing.T) {
	log.Printf("--------- Test_Null ----------")
	time.Sleep(30 * time.Second)
}

type CBFuncToTest func(*MailStation)

func createAndTest(cb CBFuncToTest) {
	log.Printf("--------- createAndTest ----------")
	srv := server.NewNode("node1")
	srv.Start(":50000")

	srv.GetCollection().Register(&test.Entry1{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Build()

	ms := NewMailStation("self")
	ms.AddServer("s0", fmt.Sprintf("localhost:%v", 50000))
	cb(ms)
	time.Sleep(3 * time.Second)
	srv.Stop()
	ms.Stop()
	time.Sleep(1 * time.Second)
}

func Test_Normal(t *testing.T) {
	createAndTest(func(ms *MailStation) {
		// service not found
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.join",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				log.Printf("rpc got result:%v %v\n", result, e)
			})))
	})
}

func Test_NoCB(t *testing.T) {
	createAndTest(func(ms *MailStation) {
		// service not found
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.join",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))})
	})
}

func Test_ServiceNotFound(t *testing.T) {
	createAndTest(func(ms *MailStation) {
		// service not found
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "abc.error",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				log.Printf("rpc got result:%v %v\n", result, e)
			})))
	})
}

func Test_MethodNotFound(t *testing.T) {
	createAndTest(func(ms *MailStation) {
		// methon not found
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.1error",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				log.Printf("rpc got result:%v %v\n", result, e)
			})))
	})
}

func Test_RPCError(t *testing.T) {
	createAndTest(func(ms *MailStation) {
		// rpc error
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.error",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				log.Printf("rpc got result:%v %v\n", result, e)
			})))
	})
}

func Test_Panic(t *testing.T) {
	createAndTest(func(ms *MailStation) {
		// rpc panic
		ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.panic",
			&test.InMsg{Abc: fmt.Sprintf("%v", rand.Intn(100))},
			common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
				log.Printf("rpc got result:%v %v\n", result, e)
			})))
	})
}

// 测试

// 服务器晚点启动
// 服务器接收到的消息顺序，一致
func Test_RPCSeq(t *testing.T) {
	log.Printf("--------- Test_RPCSeq ----------")

	srv := server.NewNode("node1")
	srv.GetCollection().Register(&test.Entry1{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Build()
	go func(){
		time.Sleep(3 * time.Second)
		srv.Start(":50000")	
	}()

	ms := NewMailStation("self")
	ms.AddServer("s0", fmt.Sprintf("localhost:%v", 50000))

	// 每秒发个消息
	index := 0

	ticker := time.NewTicker(1 * time.Millisecond)
	go func() {
		for true {
			select {
			case <-ticker.C:
				callEcho(ms, fmt.Sprintf("%v", index))
				index ++
			}
		}
	}()
	
	time.Sleep(15 * time.Second)
	srv.Stop()
	ms.Stop()
	ticker.Stop()
	time.Sleep(1 * time.Second)
}

func callEcho(ms *MailStation, str string) {	
	log.Printf("call echo:%v\n", str)
	ms.Call(fmt.Sprintf("s%v", rand.Intn(1)), "entry1.echo",
		&test.InMsg{Abc: str},
		common.ReqWithCB(reflect.ValueOf(func(result *test.OutMsg, e error) {
			log.Printf("echo got result:%v\n", result)
		})))

}


// 测试并发
// 会不会有异常
func _Test_Concurrence(t *testing.T) {
	srv := server.NewNode("node1")
	srv.Start(":50000")

	srv.GetCollection().Register(&test.Entry1{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Build()

	time.Sleep(3 * time.Second)
	ms := NewMailStation("self")
	ms.AddServer("s0", fmt.Sprintf("localhost:%v", 50000))

	time.Sleep(3 * time.Second)

	sendRPC := true

	for i := 0; i < 100; i ++ {
		go func() {
			// 随机测试
			t := time.NewTicker(1 * time.Millisecond)
			
			defer t.Stop()		

			for sendRPC {
				select {
				case <-t.C:
					testCall(ms)
				}
			}
		}()	
	}	

	time.Sleep(5 * time.Second)
	sendRPC = false

	time.Sleep(1 * time.Second)
	
	ms.Stop()		
	srv.Stop()

	time.Sleep(1 * time.Second)
	log.Printf("--------- stop ----------")
}

// 