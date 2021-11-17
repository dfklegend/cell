package tcpimpl

import (
    "encoding/json"

    "github.com/dfklegend/cell/utils/logger"
    ucommon "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/net/common/conn/message"
    //"github.com/dfklegend/cell/rpc/client"
    "github.com/dfklegend/cell/rpc/common"
    "github.com/dfklegend/cell/rpc/protos"
    "github.com/dfklegend/cell/rpc/interfaces"
    cclient "github.com/dfklegend/cell/client"
)

// grpc实现
type TCPClientImpl struct {  
    MailBox interfaces.IMailBox  
    client *cclient.CellClient  
    msgPullAcks []byte
}

func NewClientImpl() interfaces.IRPCClientImpl {
    return &TCPClientImpl{}
}

func (self *TCPClientImpl) SetMailBox(i interfaces.IMailBox) {
    //self.MailBox, _ = i.(*client.MailBox)
    self.MailBox = i
}

func (self *TCPClientImpl) Connect(address string) {
    c := cclient.NewCellClient(address)
    c.Start(address)
    self.client = c
}

func (self *TCPClientImpl) Close() {
    if self.client != nil {
        self.client.Stop()    
    }
    //self.client.Stop()    
}

func (self *TCPClientImpl) IsConnected() bool {
    return self.client != nil && self.client.IsReady()
}

func (self *TCPClientImpl) CallRPC(req *common.Request) {
    c := self.client    

    rpcReq := &protos.RPCRequest{
        ReqId:    uint32(req.ReqId),
        ClientId: self.MailBox.GetStationIId(),
        Method:   req.Method,
        Message:  ucommon.SafeJsonMarshal(req.InArg),
        NeedAck:  req.NeedAck(),
    }    

    // 发消息
    c.GetClient().SendNotify("rpc", ucommon.SafeJsonMarshalByteArray(rpcReq))
    //logger.Log.Infof("call rpc:%+v", rpcReq)
}

func (self *TCPClientImpl) ReqPullAcks() {
    if !self.IsConnected() {
        return
    }
    
    mb := self.MailBox
    c := self.client    

    //
    if len(self.msgPullAcks) == 0 {
        logger.Log.Debugf("make msgPullAcks")
        rpcReq := &protos.PullAckRequest{
            ClientId: mb.GetStationIId(),
        }    
        self.msgPullAcks = ucommon.SafeJsonMarshalByteArray(rpcReq)
    }    

    c.GetClient().SendRequest("reqpullacks", self.msgPullAcks,
        func(error bool, msg *message.Message) {
        self.onGotAcks(msg)
    })  
}

func (self *TCPClientImpl) onGotAcks(msg *message.Message) {
    r := &protos.RPCAcks{}
    err := json.Unmarshal(msg.Data, r)
    if err != nil {
        logger.Log.Debugf("onGotAcks unmarshal error:%v", err)
        return
    }

    mb := self.MailBox
    if len(r.Ack) == 0 {        
        mb.OnGotAcksCount(0)
        return       
    }  
    
    mb.OnGotAcksCount(int32(len(r.Ack)))
    mb.GetRunService().GetScheduler().Post(func() (interface{}, error) {
        self.applyAcks(r.Ack)
        return nil, nil
    })  
}

func (self *TCPClientImpl) applyAcks(acks []*protos.RPCAck) {
    mb := self.MailBox
    for _, one := range acks {
        mb.ApplyAck(one.ReqId, one.Error, one.Message)
        mb.RemoveAck(one.ReqId)
    }
}