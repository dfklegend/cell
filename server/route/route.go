// Copyright (c) nano Author and TFG Co. All Rights Reserved.
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

package route

import (
	"log"

	"github.com/dfklegend/cell/utils/logger"
	"github.com/dfklegend/cell/utils/common"
	sessioni "github.com/dfklegend/cell/server/client/session"
	"github.com/dfklegend/cell/server/interfaces"
	scommon "github.com/dfklegend/cell/server/common"
)

// 抽象，基于session或者map
type IRouteParam interface {
	Get(k string, def interface{}) interface{}
}
 
type RouteFunc func(serverType string, p IRouteParam) string

func ExamRoute(serverType string, p IRouteParam) string {
	// . 获取需求的值
	// v := p.Get("examKey").(int)	
	// . 从app获取基于服务类型的可选服务器列表
	// . 根据v来获取对应的serverId	
	return "some"
}

func GetServers(serverType string)(map[string]*scommon.Server, error) {
	app := interfaces.GetApp()
	if app == nil {
		// 模拟几个便于测试
		m := make(map[string]*scommon.Server)
		m["s1"] = &scommon.Server{
			ID: "s1",
		}
		m["s2"] = &scommon.Server{
			ID: "s2",
		}
		return m, nil
	}
	return app.GetServersByType(serverType)
}

func DefaultRoute(serverType string, p IRouteParam) string {	
	servers, err := GetServers(serverType)
	if servers == nil || len(servers) == 0 || err != nil {
		logger.Log.Errorf("can not find serverType:%v\n", serverType)
		return ""
	}
	
	var selected string
	for k, _ := range(servers) {
		selected = k
		break
	}
	return selected
}

// --------
type SessionParam struct {
	session sessioni.IServerSession
}

func NewSessionParam(session sessioni.IServerSession) *SessionParam {
	return &SessionParam {
		session: session,
	}	
}

func (self *SessionParam) Get(k string, def interface{}) interface{} {
	return self.session.Get(k, def)
}

// --------
type MapParam struct {
	data map[string]interface{}
}

func NewMapParam(data map[string]interface{}) *MapParam {
	return &MapParam {
		data: data,
	}	
}

func (self *MapParam) Get(k string, def interface{}) interface{} {
	v, ok := self.data[k]
    if ok {
        return v
    }
    return def
}

var TheRouteService = NewRouteService()
//
type IRouteService interface {
	Register(serverType string, f RouteFunc)
	GetFunc(serverType string) RouteFunc
}

type RouteService struct {
	routes map[string]RouteFunc
}

func NewRouteService() *RouteService {
	return &RouteService {
		routes: make(map[string]RouteFunc),
	}
}

func GetRouteService() *RouteService {
	return TheRouteService
}

func (self *RouteService) Register (serverType string, f RouteFunc) {
	self.routes[serverType] = f
}

func (self *RouteService) GetFunc(serverType string) RouteFunc {
	v, _ := self.routes[serverType]
	return v
}

func (self *RouteService) Route(serverType string, param interface{}) string {
	session, ok := param.(sessioni.IServerSession)
	if ok {
		return self.RouteFromSession(serverType, session)
	}

	m, ok := param.(map[string]interface{})
	if ok {
		return self.RouteFromMap(serverType, m)
	}
	log.Printf("error route param:%v\n", param)
	return ""
}

func (self *RouteService) doRoute(serverType string, param IRouteParam) string {
	f := self.GetFunc(serverType)
	if f == nil {
		f = DefaultRoute
	}

	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic in Route:%v", err)	
			log.Printf(common.GetStackStr())
		}
	}()

	return f(serverType, param)
}

func (self *RouteService) RouteFromSession(serverType string, session sessioni.IServerSession) string {
	return self.doRoute(serverType, NewSessionParam(session))	
}

func (self *RouteService) RouteFromMap(serverType string, m map[string]interface{}) string {
	return self.doRoute(serverType, NewMapParam(m))	
}