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
    "time"

    "github.com/dfklegend/cell/net/common/conn/codec"
    "github.com/dfklegend/cell/net/common/conn/message"
    "github.com/dfklegend/cell/utils/config"
)

type SessionConfig struct {
    Decoder             codec.PacketDecoder // binary decoder
    Encoder             codec.PacketEncoder // binary encoder
    MessageEncoder      message.Encoder
    HeartbeatTimeout    time.Duration
}

func NewSessionConfig(cfgFile *config.Config) *SessionConfig {
    return &SessionConfig {
        Decoder: codec.NewPomeloPacketDecoder(),
        Encoder: codec.NewPomeloPacketEncoder(),
        MessageEncoder: message.NewMessagesEncoder(cfgFile.GetBool("message.compress", false)),
    }
}

type SessionMsgConfig struct {
    MsgProcessor IMsgProcessor
}

var sessionMsgConfig = &SessionMsgConfig{}

func SetMsgProcessor(p IMsgProcessor) {
    sessionMsgConfig.MsgProcessor = p
}

func GetMsgProcessor() IMsgProcessor {
    return sessionMsgConfig.MsgProcessor
}