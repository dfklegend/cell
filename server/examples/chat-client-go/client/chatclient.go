package client

import (    
    "fmt"
    "encoding/json"
    "strings"

    "github.com/dfklegend/cell/utils/runservice" 
    "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/client"    
    "github.com/dfklegend/cell/utils/logger"
    "github.com/dfklegend/cell/net/common/conn/message"
    "github.com/dfklegend/cell/server/examples/chat-client-go/protos"
)

const (
    STATE_INIT = iota
    STATE_CONNECT_GATE
)

type ChatClient struct {
    //
    Client *client.CellClient
    RunService *runservice.StandardRunService
    state int
    msgBind bool
}

func NewChatClient(name string) *ChatClient {
    c := client.NewCellClient(name)
    r := c.GetRunService()

    return &ChatClient{
        Client: c,
        RunService: r,
        state: STATE_INIT,
        msgBind: false,
    }
}

func (self *ChatClient) setState(s int) {
    self.state = s
}

func (self *ChatClient) getState() int {
    return self.state
}

func (self *ChatClient) Start(address string) {
    c := self.Client
    c.Start(address)    

    c.SetCBConnected(self.onConnected)    
    c.SetCBBreak(func() {
        self.setState(STATE_INIT)
    })    
    self.bindMsgs()
}

func (self *ChatClient) Stop() {
    self.Client.Stop()
}

func (self *ChatClient) onConnected() {
    logger.Log.Debugf("onConnected")
    if self.getState() != STATE_INIT {
        return
    }    

    // 开始
    self.queryGate()
    self.setState(STATE_CONNECT_GATE)
}

func (self *ChatClient) bindMsgs() {
    if self.msgBind {
        return
    }
    self.msgBind = true
    // 注册
    ec := self.RunService.GetEventCenter()
    ec.Subscribe("onNewUser", func(args ...interface{}) {
        self.onNewUser(args[0].([]byte))
    })
    ec.Subscribe("onUserLeave", func(args ...interface{}) {
        //logger.Log.Debugf("onUserLeave")
        self.onUserLeave(args[0].([]byte))
    })
    ec.Subscribe("onMembers", func(args ...interface{}) {
        logger.Log.Debugf("onMembers")
    })
    ec.Subscribe("onTest", func(args ...interface{}) {
        logger.Log.Debugf("onTest")
    })
    ec.Subscribe("onMessage", func(args ...interface{}) {
        //logger.Log.Debugf("onMessage")
        self.onMessage(args[0].([]byte))
    })
}

func (self *ChatClient) queryGate() {
    logger.Log.Debugf("queryGate")

    c := self.Client
    c.GetClient().SendRequest("gate.gate.querygate", []byte("{}"), func(error bool, msg *message.Message) {
        fmt.Println("ack from cb")
        fmt.Println(string(msg.Data))

        // TODO:解析，使用正确端口
        port := self.getGatePort(msg.Data)
        c.Start(fmt.Sprintf("127.0.0.1:%v", port))
        c.GetClient().WaitReady()

        m := make(map[string]interface{})
        m["name"] = "haha"
        c.GetClient().SendRequest("gate.gate.login", []byte(common.SafeJsonMarshal(m)), self.onLoginRet);
    })
}

func (self *ChatClient) getGatePort(data []byte) string {
    obj := protos.QueryGateAck{}
    json.Unmarshal(data, &obj)    
    
    port := obj.Port
    subs := strings.Split(port, ",")
    if len(subs) >= 2 {
        return subs[1]
    }
    return ""
}

func (self *ChatClient) onLoginRet(error bool, msg *message.Message) {
    // 可以发送消息
    logger.Log.Debugf("onLoginRet")
}

func (self *ChatClient) onNewUser(data []byte) {
    obj := protos.OnNewUser{}
    json.Unmarshal(data, &obj)    
    logger.Log.Infof("onNewUser:%v", obj.Name)
}

func (self *ChatClient) onUserLeave(data []byte) {
    obj := protos.OnUserLeave{}
    json.Unmarshal(data, &obj)    
    logger.Log.Infof("onUserLeave:%v", obj.Name)
}

func (self *ChatClient) onMessage(data []byte) {
    obj := protos.ChatMsg{}
    json.Unmarshal(data, &obj)    
    logger.Log.Infof("%v: %v", obj.Name, obj.Content)
}