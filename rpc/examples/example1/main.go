package main

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"strings"
	"time"

	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/rpc/common"
	"github.com/dfklegend/cell/rpc/client"
	"github.com/dfklegend/cell/rpc/server"
)

// station1
var ms1 *client.MailStation = client.NewMailStation("station1")
var sn1 *server.RPCServerNode

func initNode1() *server.RPCServerNode {
	srv := server.NewNode("node1")
	srv.Start(":50001")

	srv.GetCollection().Register(&Service1{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Register(&Service2{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Build()
	return srv
}

func initMailStation1() {
	ms1.AddServer("node2", fmt.Sprintf("localhost:%v", 50002))
}

// station2
var ms2 *client.MailStation = client.NewMailStation("station2")
var sn2 *server.RPCServerNode

func initNode2() *server.RPCServerNode {
	srv := server.NewNode("node2")
	srv.Start(":50002")

	srv.GetCollection().Register(&Service3{}, api.WithNameFunc(strings.ToLower))
	srv.GetCollection().Build()
	return srv
}

func initMailStation2() {
	ms1.AddServer("node1", fmt.Sprintf("localhost:%v", 50001))
}

func main() {
	//
	sn1 = initNode1()
	initMailStation1()
	sn2 = initNode2()
	initMailStation2()

	var msClient *client.MailStation = client.NewMailStation("client")

	msClient.AddServer("node1", fmt.Sprintf("localhost:%v", 50001))
	msClient.AddServer("node2", fmt.Sprintf("localhost:%v", 50002))

	// call service1.func1
	msClient.Call("node1", "service1.func1",
		&InMsg{fmt.Sprintf("%v", rand.Intn(100))},
		common.ReqWithCB(reflect.ValueOf(func(result *OutMsg, e error) {
			log.Printf("client got result:%v %v\n", result, e)
		})))

	// client -> node1.service2.func1 -> node2.service3.func1 ->client
	msClient.Call("node1", "service2.func1",
		&InMsg{fmt.Sprintf("%v", rand.Intn(100))},
		common.ReqWithCB(reflect.ValueOf(func(result *OutMsg, e error) {
			// service3.func1
			log.Printf("client got result:%v %v\n", result, e)
		})))

	msClient.DumpStat()
	time.Sleep(10 * time.Second)
	msClient.DumpStat()
}
