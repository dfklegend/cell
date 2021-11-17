package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"	
	"reflect"

	"dfk.com/practices/jsongrpc/libs/util"	
	"dfk.com/practices/jsongrpc/libs/rpc"
	"dfk.com/practices/jsongrpc/cases/rpc1/protos"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)


func callRPC(ctx context.Context, client *rpc.RPCClient, method string, args interface{}) {
	log.Printf("pre call callRPC routine:%v", util.GetRoutineID())
	client.Call(ctx, method, args, reflect.ValueOf(func(err error, o *protos.OutMsg) {
		fmt.Printf("err:%v", err)
		fmt.Printf("o:%v", o)
		outMsg := o
		fmt.Printf("outMsg:%v", outMsg)
	}))	
	log.Printf("post routine:%v ", util.GetRoutineID())
}

func main() {
	
	client := &rpc.RPCClient{};
	client.Start(address);
	
	ctx, cancel := context.WithTimeout(context.Background(), 999*time.Second)
	defer cancel()
	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.GetMessage())
	//fmt.Println(name)

	for i := 0; i < 10; i++ {
		n := fmt.Sprintf("name%d", i)
		callRPC(ctx, client, "Chat.Join", &protos.InMsg{Abc:n})
	}

	for i := 0; i < 10; i++ {
		n := fmt.Sprintf("goname%d", i)
		go func(n string) {
			callRPC(ctx, client, "Chat.Join", &protos.InMsg{Abc:n})
		}(n)
	}

	exitc := make(chan os.Signal, 1)
	signal.Notify(exitc, os.Interrupt, os.Kill)

	s := <-exitc
	fmt.Println("Got signal:", s)
}
