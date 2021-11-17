package grpcimpl

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	//"github.com/dfklegend/cell/rpc/server" 
	"github.com/dfklegend/cell/rpc/protos"
	"github.com/dfklegend/cell/rpc/interfaces"
	//"github.com/dfklegend/cell/utils/common"
)

//
// 响应RPC
// 		转接到api
// 响应拉取返回结果
type GRPCServer struct {
	protos.UnimplementedRPCServerServer
	node   interfaces.IRPCNode
	server *grpc.Server
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (self *GRPCServer) Init(node interfaces.IRPCNode) {
	self.node = node
}

func (self *GRPCServer) Start(port string) {
	go self.doStart(port)
}

func (self *GRPCServer) doStart(port string) {
	log.Printf("start to listen: %v", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}
	//srv := grpc.NewServer(grpc.NumStreamWorkers(3))
	srv := grpc.NewServer()
	protos.RegisterRPCServerServer(srv, self)

	self.server = srv
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	if lis != nil {
		lis.Close()	
	}	
}

func (self *GRPCServer) Stop() {
	if self.server != nil {
		self.server.Stop()	
	}	
}

//
func (self *GRPCServer) RPC(c context.Context, in *protos.RPCRequest) (*protos.RPCReply, error) {
	//log.Printf("Route: %v Received: %v", common.GetRoutineID(), in)

	self.node.Call(in)
	// 顺便带点反馈回去
	return &protos.RPCReply{}, nil
}

func (self *GRPCServer) PullAcks(ctx context.Context, in *protos.PullAckRequest) (*protos.RPCAcks, error) {
	//log.Printf("Route: %v Receive pullAcks", common.GetRoutineID())

	acks := make([]*protos.RPCAck, 0)

	//
	rets := self.node.PopReadyAcks(in.ClientId, 400)
	if rets != nil {
		for _, wa := range rets {
			one := &protos.RPCAck{}
			one.ReqId = uint32(wa.ClientReqId)
			one.Error = wa.Err
			one.Message = wa.Message
			acks = append(acks, one)
		}
	}

	return &protos.RPCAcks{
		Ack: acks,
	}, nil
}
