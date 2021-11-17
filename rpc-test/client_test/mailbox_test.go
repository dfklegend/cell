package client

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell/utils/debug"

	"github.com/dfklegend/cell/rpc/test"
	"github.com/dfklegend/cell/rpc/client"		
)

func Test_init(t *testing.T) {	
}

// 测试gRPC断线重连(配合服务器，关服务器)
func _Test1_RunService(t *testing.T) {
	m := client.NewMailBox()
	m.Start("localhost:50051")
	fmt.Printf("%v", m)

	m.Call("Chat.Join",
		"ddd")
	ticker := TimerCall(m)

	time.Sleep(99 * time.Second)
	ticker.Stop()
	debug.DumpCurInfo()
}

func TimerCall(m *client.MailBox) *time.Ticker {
	t := time.NewTicker(1 * time.Second)
	go func() {
		for true {
			select {
			case <-t.C:
				m.Call("Chat.Join", "ddd")
			}
		}
	}()
	return t
}

func CBFunc(args *test.OutMsg, err error) (*test.OutMsg, error) {
	log.Printf("%v %v", args, err)
	return args, err
}

// 测试ack数据转化
func _Test_makeCBArgs0(t *testing.T) {
	a := &test.OutMsg{
		Def: "some data",
	}
	as, _ := json.Marshal(a)

	cbArgType := reflect.TypeOf(&test.OutMsg{})
	args := client.MakeCBArgs("", string(as), cbArgType)
	//log.Printf("%v", args[0].Interface().(*test.OutMsg))

	outMsg := args[0].Interface().(*test.OutMsg)

	//rets := reflect.ValueOf(CBFunc).Call(args)

	//log.Printf("%v", rets)
	assert.Equal(t, true, args[1].IsNil(), "空error")
	assert.Equal(t, "some data", outMsg.Def, "数据")
}

func _Test_makeCBArgs1(t *testing.T) {
	a := &test.OutMsg{
		Def: "some data",
	}
	as, _ := json.Marshal(a)

	cbArgType := reflect.TypeOf(&test.OutMsg{})
	args := client.MakeCBArgs("some error", string(as), cbArgType)
	//log.Printf("%v", args[0].Interface().(*test.OutMsg))

	//rets := reflect.ValueOf(CBFunc).Call(args)

	//log.Printf("%v", rets)
	assert.Equal(t, false, args[1].IsNil(), "空error")
}


func _TestStop(t *testing.T) {
	m := client.NewMailBox()
	m.Start("localhost:50051")
	fmt.Printf("%v", m)

	m.Call("Chat.Join",
		"ddd")
	ticker := TimerCall(m)
	time.Sleep(2 * time.Second)
	m.Stop()


	time.Sleep(5 * time.Second)
	ticker.Stop()
	debug.DumpCurInfo()
	time.Sleep(5 * time.Second)
}
