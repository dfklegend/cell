package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"dfk.com/practices/grpc/libs/util"
	pb "dfk.com/practices/grpc/protos"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func sayHallo(ctx context.Context, c pb.GreeterClient, name string) {
	log.Printf("pre call sayHallo routine:%v", util.GetRoutineID())
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("post routine:%v Greeting: %s", util.GetRoutineID(), r.GetMessage())
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), 999*time.Second)
	defer cancel()
	// r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	// log.Printf("Greeting: %s", r.GetMessage())
	fmt.Println(name)

	for i := 0; i < 10; i++ {
		n := fmt.Sprintf("name%d", i)
		sayHallo(ctx, c, n)
	}

	for i := 0; i < 10; i++ {
		n := fmt.Sprintf("goname%d", i)
		go func(n string) {
			sayHallo(ctx, c, n)
		}(n)
	}

	exitc := make(chan os.Signal, 1)
	signal.Notify(exitc, os.Interrupt, os.Kill)

	s := <-exitc
	fmt.Println("Got signal:", s)
}
