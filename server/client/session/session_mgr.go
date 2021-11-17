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

package session

import (
    "sync"
    "reflect"

    "github.com/dfklegend/cell/utils/logger"
    "github.com/dfklegend/cell/utils/common"
    nsession "github.com/dfklegend/cell/net/server/session"

    scommon "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/interfaces"
)

var frontSessions = &FrontSessionMgr{}
var SessionIdService *common.SerialIdService = common.NewSerialIdService()

func GetFrontSessions() *FrontSessionMgr {
    return frontSessions
}

func allocNetId() scommon.NetIdType {
    return scommon.NetIdType(SessionIdService.AllocId())
}

func allocCBId() scommon.NetIdType {
    return scommon.NetIdType(SessionIdService.AllocId())
}

// 管理session
// findSession
type FrontSessionMgr struct {
    sessions sync.Map
    onCloseCBs sync.Map
}

func (self *FrontSessionMgr) AddCloseCallback(cb OnSessionCloseFunc) scommon.NetIdType {    
    id := allocCBId()
    self.onCloseCBs.Store(id, cb)
    return id
}

func (self *FrontSessionMgr) DelCloseCallback(id scommon.NetIdType) {
    self.onCloseCBs.Delete(id)
}

// only test
func (self *FrontSessionMgr) findCloseCallback(cb OnSessionCloseFunc) scommon.NetIdType {

    var found scommon.NetIdType
    self.onCloseCBs.Range(func( k, v interface{}) bool {
        id := k.(scommon.NetIdType)
        c1 := v.(OnSessionCloseFunc)
        
        if reflect.ValueOf(cb).Pointer() == reflect.ValueOf(c1).Pointer() {
            found = id
            return false
        }
        return true
    })
    return found
}

func (self *FrontSessionMgr) doCloseCallback(fs *FrontSession) {
    self.onCloseCBs.Range(func( k, v interface{}) bool {
        cb := v.(OnSessionCloseFunc)
        cb(fs)
        return true
    })
}

func (self *FrontSessionMgr) AddSession(cs *nsession.ClientSession) {
    id := allocNetId()
    fs := &FrontSession {
        NetId: id,
        Session: cs,
        Data: NewSessionData(), 
    }
    self.sessions.Store(id, fs)
    cs.SetNetId(uint32(id))

    // 设置服务器id
    fs.Set(Session_ServerId, interfaces.GetApp().GetServerId())
    fs.Set(Session_NetId, id)
    logger.Log.Debugf("AddSession:%v %+v\n", id, fs)
}

func (self *FrontSessionMgr) RemoveSession(netId scommon.NetIdType) {
    s := self.FindSession(netId)
    if s != nil {        
        // 调用回调函数
        self.doCloseCallback(s)
        self.sessions.Delete(netId)
        logger.Log.Debugf("RemoveSession:%v\n", netId)    
    }
}

// 查找session
func (self *FrontSessionMgr) FindSession(netId scommon.NetIdType) *FrontSession {
    v, _ := self.sessions.Load(netId)    
    f, _ := v.(*FrontSession)
    return f
}


func AddOnSessionClose(cb OnSessionCloseFunc) {
    GetFrontSessions().AddCloseCallback(cb)
}
