package testclient

import (
	"time"
	"fmt"

	"github.com/topfreegames/pitaya/client2"
	"github.com/sirupsen/logrus"
	"github.com/topfreegames/pitaya/conn/message"
)

func NewClient() *client2.Client{
	client := client2.New(logrus.InfoLevel, 100*time.Millisecond)	
	return client;
}

// 每秒发个请求
// 处理请求
func GoTestClient() {
	client := NewClient();

	// wait connect
	err := client.ConnectTo("127.0.0.1:3250")
	if(err != nil) {
		fmt.Println("err connectto")
		return
	}
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// wait handshake
	client.WaitReady()

	// entry
	client.SendRequest("room.room.entry", []byte("h1"), func(error bool, msg *message.Message) {
		fmt.Println("ack from cb")
		fmt.Println(string(msg.Data))
	})
	
	for true {
		select {
		case <-ticker.C:
			client.SendRequest("room.room.hello", []byte("h1"), func(error bool, msg *message.Message) {
				fmt.Println("ack from cb")
				fmt.Println(string(msg.Data))
			})
		case msg :=<- client.MsgChannel():
			if(msg.Cb != nil) {
				msg.Cb(false, msg.Msg)
			}
		}

	}
}

func TestClients(num int) {
	for i := 0; i < num; i ++ {
		go GoTestClient();
	}
}

