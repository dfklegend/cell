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
    "encoding/json"
)

// 方便打包传到后端服务器构建backsession
type SessionData struct {
    data map[string]interface{}
}

func NewSessionData() *SessionData {
    return &SessionData {
        data: make(map[string]interface{}),
    }
}

func (self *SessionData) GetMap() map[string]interface{} {
    return self.data
}

func (self *SessionData) Set(k string, v interface{}) {
    self.data[k] = v
}

func (self *SessionData) Get(k string, def interface{}) interface{} {
    v, ok := self.data[k]
    if ok {
        return v
    }
    return def
}

func (self *SessionData) Has(k string) bool {
    _, ok := self.data[k]
    return ok
}

func (self *SessionData) ToJson() []byte {
    v, _ := json.Marshal(self.data)
    return v
}

func (self *SessionData) FromJson(d []byte) {
    json.Unmarshal(d, &self.data)
}

func (self *SessionData) ToJsonStr() string {
    return string(self.ToJson())
}

func (self *SessionData) FromJsonStr(s string) {
    json.Unmarshal([]byte(s), &self.data)
}