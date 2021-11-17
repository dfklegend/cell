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
    "encoding/json"
    
    nsession "github.com/dfklegend/cell/net/server/session"
    "github.com/dfklegend/cell/server/common"
    "github.com/dfklegend/cell/server/interfaces"
)

// frontsession, backsession，为backend服务器提供

// 可以直接通过session，向客户端推送消息
// 向客户端推送消息，已知[serverId,netId]即可
type IServerSession interface {
    // 绑定ID
    Bind(id string)
    GetID() string
    // 设置值
    Set(k string,v interface{})
    Get(k string, def interface{}) interface{}
    // 推送更新
    PushSession()
    ToJson() string
    FromJson(string)

    // 踢下线
    Kick()

    // 互斥锁    
    // 使用Bind,GetID,Set/Get
    // ToJson,FromJson
    Lock()
    Unlock()
}

// 回调函数
type OnSessionCloseFunc func(s IServerSession)

const (
    Session_ID = "_ID"
    Session_NetId = "_NetId"
    Session_ServerId = "_ServerId"
)

// ----------------------
// 保留
//      _ID         Bind的唯一ID
//      _NetId      对应的frontSession id
//      _ServerId   对应的frontserver id
type FrontSession struct {
    NetId       common.NetIdType
    Session     *nsession.ClientSession
    Data        *SessionData 
    mutex       sync.Mutex
}

func NewFrontSession() *FrontSession {
    return &FrontSession{
        Data: NewSessionData(),
    }
}

func (self *FrontSession) Reserve() {}
func (self *FrontSession) Handle() {}

func (self *FrontSession) Bind(id string) {
    self.Data.Set(Session_ID, id)
}

func (self *FrontSession) GetID() string {
    return self.Get(Session_ID, "").(string)
}

func (self *FrontSession) Set(k string, v interface{}) {
    self.Data.Set(k, v)
}

func (self *FrontSession) Get(k string, def interface{}) interface{} {
    return self.Data.Get(k, def)
}

func (self *FrontSession) PushSession() {
    // doNothing
}

func (self *FrontSession) ToJson() string {
    return self.Data.ToJsonStr()
}

func (self *FrontSession) FromJson(d string) {
    self.Data.FromJsonStr(d)
}

func (self *FrontSession) Kick() {
    self.Session.Close()
}

func (self *FrontSession) Lock() {
    self.mutex.Lock()
}

func (self *FrontSession) Unlock() {
    self.mutex.Unlock()
}

// ----------------------
// 需要知道session所在serverId和NetId
type BackSession struct {    
    ServerId    string
    NetId       common.NetIdType
    dirt        bool
    Data        *SessionData 
    NewData        *SessionData
}

func NewBackSession(str string) *BackSession {
    v := &BackSession{
        Data: NewSessionData(),
        NewData: NewSessionData(),
    }
    v.FromJson(str)
    return v
}

// 可以直接创建一个backSession，用于push session
func NewBackSessionForPush(serverId string, netId common.NetIdType) *BackSession {
    v := &BackSession{
        Data: NewSessionData(),
        NewData: NewSessionData(),
    }
    v.ServerId = serverId
    v.NetId = netId
    return v
}

func (self *BackSession) Reserve() {}
func (self *BackSession) Handle() {}

func (self *BackSession) Bind(id string) {
    self.Set(Session_ID, id)
}

func (self *BackSession) GetID() string {
    return self.Get(Session_ID, "").(string)
}

func (self *BackSession) Set(k string, v interface{}) {
    self.NewData.Set(k, v)
    self.dirt = true
}

func (self *BackSession) Get(k string, def interface{}) interface{} {
    if self.NewData.Has(k) {
        return self.NewData.Get(k, def)
    }
    return self.Data.Get(k, def)
}

// 避免推送了过多的数据
func (self *BackSession) PushSession() {
    // 推送给目标
    if !self.dirt {
        return
    }
    self.dirt = false
    ss := interfaces.GetApp().GetSysService()
    ss.ReqPushSession(self.ServerId, self.NetId, self.NewData.ToJsonStr(), nil)
}

func (self *BackSession) ToJson() string {
    // 合成一个全的
    m1 := self.Data.GetMap()
    m2 := self.NewData.GetMap()

    newMap := make(map[string]interface{})
    mergeMap(newMap, m1)
    mergeMap(newMap, m2)

    v, _ := json.Marshal(newMap)
    return string(v)
}

func mergeMap(m1, m2 map[string]interface{}) {
    for k, v := range m2 {
        m1[k] = v
    }
}

func (self *BackSession) FromJson(d string) {
    self.Data.FromJsonStr(d)
    self.ServerId = self.Get(Session_ServerId, "n").(string)
    if self.Data.Has(Session_NetId) {
        // json转interface value时，number被转成float64
        self.NetId = common.NetIdType(self.Get(Session_NetId, 0).(float64))    
    }
    
    self.dirt = false
}

func (self *BackSession) Kick() {
    ss := interfaces.GetApp().GetSysService()
    ss.ReqKick(self.ServerId, self.NetId, nil)
}

func (self *BackSession) Lock() {
}

func (self *BackSession) Unlock() {
}


// helper funcs
func GetServerId(s IServerSession) string {
    return s.Get(Session_ServerId, "").(string)
}

func GetNetId(s IServerSession) common.NetIdType {
    return s.Get(Session_NetId, 0).(common.NetIdType)
}



