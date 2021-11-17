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

package services

import (    
    "reflect"    
    "errors"

    "github.com/dfklegend/cell/utils/logger"
    rpccommon "github.com/dfklegend/cell/rpc/common"
    api "github.com/dfklegend/cell/rpc/apientry"
    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/net/common/conn/message"
    ninterfaces "github.com/dfklegend/cell/net/interfaces"
    pkgsession "github.com/dfklegend/cell/server/client/session"
    "github.com/dfklegend/cell/server/interfaces"
    "github.com/dfklegend/cell/server/services/sys"
    "github.com/dfklegend/cell/server/route"
)

var (
    emptyResult  = []byte("{}")
    routeErrorNoServer = errors.New("route error, no server to route")
)


// 负责根据转发给本地接口或者远程接口
// 负责远程接口返回的对接
type MsgProcessor struct {
    //
}

func (self *MsgProcessor) ProcessMessage(s ninterfaces.IHandlerSession, msg *message.Message) {
    fs := s.(*pkgsession.FrontSession)    
    app := interfaces.GetApp()
    routeStr := common.NewRouteStr(msg.Route)
    // 判断是否同一类型
    logger.Log.Debugf("routeStr:%v", routeStr)
    if app.GetServerType() == routeStr.Server {
        self.localMessage(fs, msg)
        return    
    }

    clientSession := fs.Session    
    // 路由
    selected := route.GetRouteService().RouteFromSession(routeStr.Server, fs)
    if selected == "" {
        // 路由失败
        if msg.ID > 0 {
            clientSession.ResponseMID(msg.ID, emptyResult, routeErrorNoServer)    
        }
        return
    }

    logger.Log.Errorf("route selected:%v", selected)

    ms := app.GetMailStation()    
    // 调用远程接口
    ms.Call(selected, "sys.call", &sys.SysReq{        
            SessionData: fs.ToJson(),
            Method: routeStr.GetShortStr(),
            Arg: string(msg.Data),
        }, rpccommon.ReqWithCB(reflect.ValueOf(func(result *sys.SysAck, e error) {

        if e != nil {
            logger.Log.Errorf("sys.call error:%v", e)
        }
        if msg.ID > 0 {
            clientSession.ResponseMID(msg.ID, []byte(result.Result), e)    
        }   
    })))
}

func (self *MsgProcessor) localMessage(s api.IHandlerSession, msg *message.Message) {
    fs := s.(*pkgsession.FrontSession)
    session := fs.Session

    routeStr := common.NewRouteStr(msg.Route)
    // 判断是否同一类型
    logger.Log.Debugf("routeStr:%v", routeStr)
    // 要call handler
    interfaces.GetApp().GetHandler().Call(fs,
        routeStr.GetShortStr(), msg.Data, func(e error, result interface{}) {
        logger.Log.Debugf("handler.call:%v cb result:%v", routeStr, result)
        // 返回值
        if msg.ID > 0 {
            session.ResponseMID(msg.ID, result, e)    
        }        
    })
}

