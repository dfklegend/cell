package tcpimpl

import (
    "time"
    "reflect"
    "encoding/json"

    "github.com/dfklegend/cell/utils/logger"
    ucommon "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/net/server/acceptor"
    "github.com/dfklegend/cell/net/common/conn/message"
    "github.com/dfklegend/cell/net/interfaces"
    nsession "github.com/dfklegend/cell/net/server/session"
    //"github.com/dfklegend/cell/rpc/server"
    "github.com/dfklegend/cell/rpc/protos"
    rpcinterfaces "github.com/dfklegend/cell/rpc/interfaces"
)

// 启动一个tcp监听者
// 处理
type TCPServerImpl struct {  
    node rpcinterfaces.IRPCNode
    server acceptor.Acceptor
    sessionCfg *nsession.SessionConfig
}

func NewServerImpl() rpcinterfaces.IRPCServerImpl {
    return &TCPServerImpl{}
}

func (self *TCPServerImpl) Init(node rpcinterfaces.IRPCNode) {
    //n, _ := node.(*server.RPCServerNode)
    self.node = node
}

func (self *TCPServerImpl) Start(listenAddress string) {

    self.sessionCfg = nsession.NewSessionConfig(nil)
    self.sessionCfg.Impl = NewTCPServerSessionImpl(self.node)

    self.server = acceptor.NewTCPAcceptor(listenAddress)
    self.startAcceptorListen(self.server)
}

func (self *TCPServerImpl) Stop() {
    self.server.Stop()    
}

func (self *TCPServerImpl) startAcceptorListen(a acceptor.Acceptor) {
    // new session
    go func() {
        for conn := range a.GetConnChan() {
            logger.Log.Debugf("new conn come:%v", conn)
            // 新连接建立                
            s := nsession.NewClientSession(conn,
                self.sessionCfg)
            // add to 
            s.Handle()
            logger.Log.Debugf("%v", conn)
        }
    }()

    go func() {
        // 监听
        a.ListenAndServe()
    }()

    go func() {
        time.Sleep(time.Second)
        logger.Log.Infof("listening with acceptor %s on addr %s", reflect.TypeOf(a), a.GetAddr())   
    }()     
}

// ---------------
type TCPServerSessionImpl struct {
    node rpcinterfaces.IRPCNode
}

func NewTCPServerSessionImpl(n rpcinterfaces.IRPCNode) *TCPServerSessionImpl {
    return &TCPServerSessionImpl{
        node: n,
    }
}

func (self *TCPServerSessionImpl) ProcessMessage(s interfaces.IClientSession, msg *message.Message) {
    // 处理消息
    //logger.Log.Infof("ProcessMessage:%+v", msg)
    // rpc
    if msg.Route == "rpc" {
        self.onRPC(s, msg)
        return
    }
    if msg.Route == "reqpullacks" {
        self.onPullAcks(s, msg)
        return
    }    
}
    
func (self *TCPServerSessionImpl) OnSessionCreate(s interfaces.IClientSession) {

}

func (self *TCPServerSessionImpl) OnSessionClose(s interfaces.IClientSession) {

}

func (self *TCPServerSessionImpl) onRPC(s interfaces.IClientSession, msg *message.Message) {
    req := &protos.RPCRequest{}
    err := json.Unmarshal(msg.Data, req)
    if err != nil {
        logger.Log.Warnf("error unmarshal")
        return
    }

    //logger.Log.Infof("onRPC req:", req)
    self.node.Call(req)
}

func (self *TCPServerSessionImpl) onPullAcks(i interfaces.IClientSession, msg *message.Message) {
    cs, _ := i.(*nsession.ClientSession)
    if cs == nil {
        return
    }

    req := &protos.PullAckRequest{}
    err := json.Unmarshal(msg.Data, req)
    if err != nil {
        logger.Log.Warnf("error unmarshal")
        return
    }    

    acks := make([]*protos.RPCAck, 0)
    rets := self.node.PopReadyAcks(req.ClientId, 400)
    if rets != nil {
        for _, wa := range rets {
            one := &protos.RPCAck{}
            one.ReqId = uint32(wa.ClientReqId)
            one.Error = wa.Err
            one.Message = wa.Message
            acks = append(acks, one)
        }
    }

    ackMsg := &protos.RPCAcks{
        Ack: acks,
    }

    // 序列化,并返回
    cs.ResponseMID(msg.ID, ucommon.SafeJsonMarshalByteArray(ackMsg), nil)
}