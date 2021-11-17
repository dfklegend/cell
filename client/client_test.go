package client

import (
	"testing"
	"time"	
	"fmt"

	"github.com/dfklegend/cell/net/common/conn/message"
)

// 功能测试
func TestCellClient(t *testing.T) {
	c := NewCellClient("client")
	go func() {
		c.Start("127.0.0.1:30021")
		c.GetClient().WaitReady()

		c.GetClient().SendRequest("gate.gate.querygate", []byte("{}"), func(error bool, msg *message.Message) {
			fmt.Println("ack from cb")
			fmt.Println(string(msg.Data))
		})
	}()
	
	time.Sleep(5*time.Second)
	c.Stop()
}
