package rpc

import (
	"context"	
	"log"
	"fmt"
	"encoding/json"
	"reflect"
	"errors"

	"google.golang.org/grpc"

	"dfk.com/practices/jsongrpc/libs/util"
	pb "dfk.com/practices/jsongrpc/protos"
)

// 是否错误, 返回类型
type CBFunc reflect.Value


type RPCClient struct {
	conn *grpc.ClientConn
	rpcClient pb.RPCServerClient
}

func (self *RPCClient) Start(address string) {
	c, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	self.conn = c;

	c1 := pb.NewRPCServerClient(c)
	self.rpcClient = c1
}

func (self *RPCClient) Call(ctx context.Context, method string,
	inArgs interface{}, cb reflect.Value) {
	if self.conn == nil {
		return
	}
	
	data, err := json.Marshal(inArgs)
	r, err := callRPC(ctx, self.rpcClient, method, string(data))

	//cbFunc := cb.Interface()

	// 根据cbFunc的参数来反射，先直接调回去
	if err != nil {
		//cbFunc(err, nil)
		return
	}

	// 定义错误
	if r.Code != 0 {
		//cbFunc(nil, nil)
		return
	}

	//ret := reflect.New(cbArgType.Elem()).Interface()
	//err = json.Unmarshal([]byte(r.Message), ret)

	//cb(err, ret)
	outData, _ := tryParseFuncArgs(r.Message, cb)
	
	//cbFunc(nil, data)
	// https://stackoverflow.com/questions/26321115/using-reflection-to-call-a-function-with-a-nil-parameter-results-in-a-call-usin
	//var a error
	//nilValue := reflect.New(reflect.TypeOf(a)).Elem()   // error~~
	nilValue := reflect.New(reflect.TypeOf(errors.New(""))).Elem()
	//args := []reflect.Value{reflect.ValueOf(errors.New("")), reflect.ValueOf(outData)}
	args := []reflect.Value{nilValue, reflect.ValueOf(outData)}
	log.Printf("args:%v", args)
	cb.Call(args)
}

func callCB(cb reflect.Value) {

}

func callRPC(ctx context.Context, c pb.RPCServerClient,
	 method string, message string)(*pb.RPCReply, error) {
	log.Printf("pre call callRPC routine:%v", util.GetRoutineID())
	r, err := c.RPC(ctx, &pb.RPCRequest{Method: method, Message:message})
	if err != nil {
		log.Fatalf("could not call rpc: %v", err)
	}
	log.Printf("post routine:%v Greeting: %v", util.GetRoutineID(), r)
	return r, err
}

func tryParseFuncArgs(inData string, inCb reflect.Value)(interface{}, error){
	fmt.Printf("inCb: %v\n", inCb)
	cb := inCb.Interface()
	fmt.Printf("cb: %v\n", cb)
	typ := reflect.TypeOf(cb);
	fmt.Printf("type:%v\n", typ)
	if typ.NumIn() != 2 {
		return "", nil
	}	

	argType := typ.In(1);
	data := reflect.New(argType.Elem()).Interface()
	err := json.Unmarshal([]byte(inData), data)
	return data, err
}