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

package interfaces

import (    
    "github.com/dfklegend/cell/net/common/conn/message"     
)

// = 
type IHandlerSession interface {
    Reserve()
    Handle()
}

type IMsgProcessor interface{    
    ProcessMessage(IHandlerSession, *message.Message)
}


// 代表clientSession
type IClientSession interface {
    Reserve()
    Push(route string, v interface{}) error
}

// 外部传入，处理session消息
type IClientSessionImpl interface {
    // 处理接收的消息
    ProcessMessage(IClientSession, *message.Message)
    OnSessionCreate(IClientSession)
    OnSessionClose(IClientSession)
}