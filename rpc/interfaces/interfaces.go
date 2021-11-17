package interfaces

import (
    "github.com/dfklegend/cell/utils/runservice"

    "github.com/dfklegend/cell/rpc/common"
    "github.com/dfklegend/cell/rpc/protos"
)

type IMailBox interface {    
    GetStationIId() string
    ApplyAck(reqId uint32, errStr string, message string)
    RemoveAck(reqId uint32)
    OnGotAcksCount(acksCount int32)
    GetRunService() *runservice.RunService
}

type IRPCClientImpl interface {
    SetMailBox(IMailBox)
    Connect(address string)
    Close()
    IsConnected() bool
    CallRPC(req *common.Request)
    ReqPullAcks()    
}

type IRPCNode interface {   
    Call(in *protos.RPCRequest) 
    PopReadyAcks(clientId string, maxNum int) []*common.WaitAckReq
}

type IRPCServerImpl interface {
    Init(node IRPCNode)
    Start(listenAddress string)
    Stop()
}
