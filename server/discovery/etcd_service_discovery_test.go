// Copyright (c) TFG Co. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package discovery

import (
	"testing"
	"time"

	"github.com/dfklegend/cell/server/common"
)

var etcdSDTables = []struct {
	server *common.Server
}{
	{common.NewServer("frontend-1", "type1", "", map[string]string{"k1": "v1"})},
	{common.NewServer("backend-1", "type2", "", map[string]string{"k2": "v2"})},
	{common.NewServer("backend-2", "type3", "", nil)},
}

var etcdSDTablesMultipleServers = []struct {
	servers []*common.Server
}{
	{[]*common.Server{}},
	{[]*common.Server{
		common.NewServer("frontend-1", "type1", "", map[string]string{"k1": "v1"}),
		common.NewServer("backend-1", "type2", "", map[string]string{"k2": "v2"}),
		common.NewServer("backend-2", "type3", "", nil),
	}},
	{[]*common.Server{
		common.NewServer("frontend-1", "type2", "", map[string]string{"k1": "v1"}),
		common.NewServer("frontend-2", "type2", "", map[string]string{"k2": "v2"}),
		common.NewServer("frontend-3", "type2", "", map[string]string{"k1": "v1"}),
	}},
}

var etcdSDBlacklistTables = []struct {
	name                string
	server              *common.Server
	serversToAdd        []*common.Server
	serverTypeBlacklist []string
}{
	{
		name:   "test1",
		server: common.NewServer("frontend-1", "type1", "", nil),
		serversToAdd: []*common.Server{
			common.NewServer("frontend-1", "type1", "", nil),
		},
		serverTypeBlacklist: nil,
	},
	{
		name:   "test2",
		server: common.NewServer("frontend-1", "type1", "", nil),
		serversToAdd: []*common.Server{
			common.NewServer("backend-1", "type1", "", nil),
			common.NewServer("backend-2", "type2", "", nil),
			common.NewServer("backend-3", "type3", "", nil),
		},
		serverTypeBlacklist: []string{"type2"},
	},
	{
		name:   "test3",
		server: common.NewServer("frontend-1", "type1", "", nil),
		serversToAdd: []*common.Server{
			common.NewServer("backend-1", "type1", "", nil),
			common.NewServer("backend-2", "type2", "", nil),
			common.NewServer("backend-3", "type3", "", nil),
			common.NewServer("backend-4", "type4", "", nil),
		},
		serverTypeBlacklist: []string{"type1", "type4"},
	},
}

func TestEtcdInit(t *testing.T) {
	config := NewDefaultEtcdServiceDiscoveryConfig()
	server := etcdSDTables[0].server

	appDieChan := make(chan bool)
	e, _ := NewEtcdServiceDiscovery(*config, server, appDieChan, nil)

	etcd := e.(*ETCDServiceDiscovery)
	etcd.Init()

	time.Sleep(10 * time.Second)
}

func runETCD(server *common.Server) {
	config := NewDefaultEtcdServiceDiscoveryConfig()
	appDieChan := make(chan bool)
	e, _ := NewEtcdServiceDiscovery(*config, server, appDieChan, nil)

	etcd := e.(*ETCDServiceDiscovery)
	etcd.Init()

	time.Sleep(10 * time.Second)
}

func TestMulti(t *testing.T) {
	t.Parallel()
	for _, table := range etcdSDTables {
		t.Run(table.server.ID, func(t *testing.T) {
			go runETCD(table.server)
		})
	}

	time.Sleep(20 * time.Second)
}
