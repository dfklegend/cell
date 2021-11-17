/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	
	"dfk.com/practices/jsongrpc/libs/rpc"
	"dfk.com/practices/jsongrpc/libs/rpc/service"	
	"dfk.com/practices/jsongrpc/cases/rpc1/server/services"
)

const (
	port = ":50051"
)



func main() {
	

	service.NewMgr()
	service.TheMgr.Register(&services.Chat{})
	service.TheMgr.Register(&services.Some{})
	service.TheMgr.Start()

	server := &rpc.RPCServer{}
	server.Start(port);	
}
