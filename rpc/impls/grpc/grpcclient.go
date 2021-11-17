package grpcimpl

import (
    "context"
    "encoding/json"

    "google.golang.org/grpc"
    "google.golang.org/grpc/connectivity"

    "github.com/dfklegend/cell/utils/logger"

    //"github.com/dfklegend/cell/rpc/client"
    "github.com/dfklegend/cell/rpc/common"
    "github.com/dfklegend/cell/rpc/protos"
    "github.com/dfklegend/cell/rpc/interfaces"
)

// grpc实现
type GRPCClientImpl struct {  
    MailBox interfaces.IMailBox
    conn *grpc.ClientConn  
    rpcClient protos.RPCServerClient
}

func NewClientImpl() interfaces.IRPCClientImpl {
    return &GRPCClientImpl{}
}

func (self *GRPCClientImpl) SetMailBox(i interfaces.IMailBox) {
    self.MailBox = i
}

func (self *GRPCClientImpl) Connect(address string) {
    logger.Log.Infof("%v mailbox start connect\n", address)
    //c, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    c, err := grpc.DialContext(
        context.Background(), address, grpc.WithInsecure())
    logger.Log.Infof("%v mailbox connect over\n", address)
    if err != nil {
        logger.Log.Errorf("did not connect: %v", err)
    }
    self.conn = c

    c1 := protos.NewRPCServerClient(c)
    self.rpcClient = c1
}

func (self *GRPCClientImpl) Close() {
    if self.conn == nil {
        return
    }   
    self.conn.Close()
}

func (self *GRPCClientImpl) IsConnected() bool {
    return self.conn != nil && self.conn.GetState() == connectivity.Ready
}

func (self *GRPCClientImpl) CallRPC(req *common.Request) {
    //logger.Log.Infof("State:%v request:%v", self.conn.GetState().String(),
    //    req.ReqId)

    c := self.rpcClient
    data, err := json.Marshal(req.InArg)
    _, err = c.RPC(context.Background(),
        &protos.RPCRequest{
            ReqId:    uint32(req.ReqId),
            ClientId: self.MailBox.GetStationIId(),
            Method:   req.Method,
            Message:  string(data),
            NeedAck:  req.NeedAck(),
        })
    if err != nil {
        // TODO: 重试?
        logger.Log.Warnf("could not call rpc: %v", err)
    }
}

func (self *GRPCClientImpl) ReqPullAcks() {
    if !self.IsConnected() {
        return
    }

    mb := self.MailBox
    c := self.rpcClient
    //
    r, _ := c.PullAcks(context.Background(),
        &protos.PullAckRequest{
            ClientId: mb.GetStationIId(),
        })

    if r == nil || len(r.Ack) == 0 {        
        mb.OnGotAcksCount(0)
        return
    }    

    mb.OnGotAcksCount(int32(len(r.Ack)))
    mb.GetRunService().GetScheduler().Post(func() (interface{}, error) {
        self.applyAcks(r.Ack)
        return nil, nil
    })      
}

func (self *GRPCClientImpl) applyAcks(acks []*protos.RPCAck) {
    mb := self.MailBox
    for _, one := range acks {
        mb.ApplyAck(one.ReqId, one.Error, one.Message)
        mb.RemoveAck(one.ReqId)
    }
}