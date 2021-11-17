package server

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"testing"
	"time"

	"github.com/dfklegend/cell/rpc/client"
	"github.com/dfklegend/cell/rpc/consts"
	"github.com/dfklegend/cell/rpc/protos"	
	"github.com/dfklegend/cell/utils/common"
	"google.golang.org/grpc"
)

// 用来调试客户端
type TestRPCServer struct {
	protos.UnimplementedRPCServerServer
	server *grpc.Server
}

func (s *TestRPCServer) Start(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer(grpc.NumStreamWorkers(3))
	s.server = srv
	protos.RegisterRPCServerServer(srv, s)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *TestRPCServer) Stop() {
	if s.server != nil {
		s.server.Stop()	
	}	
}

//
func (s *TestRPCServer) RPC(c context.Context, in *protos.RPCRequest) (*protos.RPCReply, error) {
	log.Printf("Route: %v Received: %v", common.GetRoutineID(), in.GetMessage())

	return &protos.RPCReply{}, nil
}

func (c *TestRPCServer) PullAcks(ctx context.Context, in *protos.PullAckRequest) (*protos.RPCAcks, error) {
	//log.Printf("Route: %v Receive pullAcks", common.GetRoutineID())

	acks := make([]*protos.RPCAck, 0)
	//acks = append(acks, newRPCAck())
	//acks = append(acks, newRPCAck())
	return &protos.RPCAcks{
		Ack: acks,
	}, nil
}

func newRPCAck() *protos.RPCAck {
	return &protos.RPCAck{
		ReqId:   1,
		Message: "haha",
	}
}

func FuncTest_TestServer(t *testing.T) {
	log.Printf("Test server start!\n")
	server := &TestRPCServer{}
	server.Start(consts.DefaultServerPort)
}

func StartTestServer(port string) *TestRPCServer {
	srv := &TestRPCServer{}
	srv.Start(port)
	return srv
}

// 测试多个目标
func Test_mailStation(t *testing.T) {
	ms := client.NewMailStation("node1")

	servers := make([]*TestRPCServer,0)
	chanServer := make (chan *TestRPCServer, 99)

	for i := 0; i < 5; i++ {
		go func(index int) {
			s := StartTestServer(fmt.Sprintf(":%v", 50000+index))
			chanServer <- s
		}(i)
		
		ms.AddServer(fmt.Sprintf("s%v", i), fmt.Sprintf("localhost:%v", 50000+i))
	}

	go func(){
		for s:= range chanServer {
			servers = append(servers, s)
		}
	}()



	ms.Call("s0", "entry1.join", "ddd")

	sendRPC := true
	go func() {
		// 随机测试
		t := time.NewTicker(1 * time.Millisecond)

		for sendRPC {
			select {
			case <-t.C:
				ms.Call(fmt.Sprintf("s%v", rand.Intn(6)), "entry1.join", "ddd")
			}
		}
	}()

	time.Sleep(10 * time.Second)
	sendRPC = false

	for _, s := range(servers) {
		s.Stop()
	}

	time.Sleep(1 * time.Second)
}
