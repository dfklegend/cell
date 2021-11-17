package cmds

import (
	"fmt"

	pc "github.com/dfklegend/cell/cmdclient/pitayaclient"
	"github.com/dfklegend/cell/cmdclient/testclient"
	"github.com/topfreegames/pitaya/conn/message"
)

func init() {
	Visit()
	RegisterCmd(&ConnectCmd{})
	RegisterCmd(&RequestCmd{})
	RegisterCmd(&ClientsCmd{})
}

type ConnectCmd struct {
}

func (s *ConnectCmd) GetName() string {
	return "connect"
}

func (s *ConnectCmd) Do(args []string) {
	pc.Start()
	pc.GetClient().ConnectTo("127.0.0.1:3250")
}

type RequestCmd struct {
}

func (s *RequestCmd) GetName() string {
	return "request"
}

func (s *RequestCmd) Do(args []string) {
	
	pc.GetClient().SendRequest("room.room.entry", []byte("h1"), func(error bool, msg *message.Message) {
		fmt.Println("ack from cb")
		fmt.Println(string(msg.Data))
	})
	pc.GetClient().SendRequest("room.room.entry", []byte("h2"), func(error bool, msg *message.Message) {
		fmt.Println("ack1 from cb")
		fmt.Println(string(msg.Data))
	})

	go s.getMsgs()
}

func (s *RequestCmd) getMsgs() {
	for msg:= range pc.GetClient().MsgChannel() {
		//fmt.Printf("msg:%v", msg)
		if(msg.Cb != nil) {
			msg.Cb(false, msg.Msg)
		}

	}
}


type ClientsCmd struct {
}

func (s *ClientsCmd) GetName() string {
	return "clients"
}

func (s *ClientsCmd) Do(args []string) {
	testclient.TestClients(99)
}

