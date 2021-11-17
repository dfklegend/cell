
package rpc

import (
	"context"
	"log"
	"net"	
	"fmt"	
	"encoding/json"

	"google.golang.org/grpc"

	"dfk.com/practices/jsongrpc/libs/util"
	"dfk.com/practices/jsongrpc/libs/rpc/service"
	pb "dfk.com/practices/jsongrpc/protos"
	"dfk.com/practices/jsongrpc/cases/rpc1/protos"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type RPCServer struct {
	pb.UnimplementedRPCServerServer	
}

func (s *RPCServer) Start(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.NumStreamWorkers(3))
	pb.RegisterRPCServerServer(srv, s)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// 
func (s *RPCServer) RPC(c context.Context, in *pb.RPCRequest) (*pb.RPCReply, error) {
	log.Printf("Route: %v Received: %v", util.GetRoutineID(), in.GetMessage())

	// 投递到目标service,做参数匹配	 		
	var ret = protos.OutMsg{}
	ret.Def = "hahah";

	data, _:= json.Marshal(ret)

	err := service.TheMgr.Call(in.Method, []byte(in.Message) )
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	return &pb.RPCReply{Code: 0, Message: string(data)}, nil
}
