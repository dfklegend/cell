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
    "fmt"
    "time"
    "sync"
    "sync/atomic"
    "encoding/json"
    "errors"

    "github.com/dfklegend/cell/utils/logger"
    "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/utils/compression"
    
    "github.com/dfklegend/cell/net/common/conn/codec"
    "github.com/dfklegend/cell/net/common/conn/packet"
    "github.com/dfklegend/cell/net/common/conn/message"    
    "github.com/dfklegend/cell/net/constants"
    "github.com/dfklegend/cell/net/interfaces"
    "github.com/dfklegend/cell/net/server/acceptor" 
    //"github.com/dfklegend/cell/server/interfaces"    
)

// session state
const (
    _ int32 = iota
    // StatusStart status
    StatusStart
    // StatusHandshake status
    StatusHandshake
    // StatusWorking status
    StatusWorking
    // StatusClosed status
    StatusClosed
)

var (
    // hbd contains the heartbeat packet data
    hbd []byte
    // hrd contains the handshake response data
    hrd  []byte
    once sync.Once

    DefaultHeartbeatTimeSeconds int64 = 10
    DefaultHeartbeatTime = 10 * time.Second
)

type pendingWrite struct {    
    data []byte
    err  error
}

type pendingMessage struct {    
    typ     message.Type // message type
    route   string       // message route (push)
    mid     uint         // response message id (response)
    payload interface{}  // payload
    err     bool         // if its an error message
}

//
type ISession interface{
    Handle()
}


// 负责接收来自客户端的消息
type ClientSession struct {
    // set after add to frontSessions
    netId           uint32
    // 连接对象
    conn            acceptor.PlayerConn
    // 发送队列
    chSend          chan *pendingWrite

    cfg             *SessionConfig
    // 注:目前是nil
    // 因为下发数据都是已经序列化过了[]byte
    //serializer      serialize.Serializer

    chanClose       chan bool
    mutex           sync.Mutex

    state           int32                // current agent state
    lastHeartBeat   int64
}

func NewClientSession(        
        c acceptor.PlayerConn,
        cfg *SessionConfig) *ClientSession {

    once.Do(func() {
        serializerName := "json"
        heartbeatTime := DefaultHeartbeatTime
        hbdEncode(heartbeatTime, cfg.Encoder, 
            cfg.MessageEncoder.IsCompressionEnabled(), serializerName)
    })

    v := &ClientSession{
        conn: c,
        cfg: cfg,        
        chSend: make(chan *pendingWrite, 9999), 
        chanClose: make(chan bool), 
        state: StatusStart,
    }   

    // TODO
    //GetFrontSessions().AddSession(v)
    v.GetImpl().OnSessionCreate(v)
    return v 
}

func (self *ClientSession) GetImpl() interfaces.IClientSessionImpl {
    return self.cfg.Impl
}

func (self *ClientSession) Reserve() {}

func (self *ClientSession) SetNetId(id uint32) {
    self.netId = id
}

func (self *ClientSession) GetNetId() uint32 {
    return self.netId
}

func (self *ClientSession) SetStatus(state int32) {
    atomic.StoreInt32(&self.state, state)
}

// GetStatus gets the status
func (self *ClientSession) GetStatus() int32 {
    return atomic.LoadInt32(&self.state)
}

func (self *ClientSession) Close() { 
    self.mutex.Lock()
    defer self.mutex.Unlock()

    select {
    // close already
    case <- self.chanClose:
        return
    default:
        close(self.chanClose)
        close(self.chSend)
    }
    // close channel
    self.conn.Close()   

    self.GetImpl().OnSessionClose(self)
    //GetFrontSessions().RemoveSession(self.netId)
}


// 读写消息
func (self *ClientSession) Handle() {    
    go self.heartbeat();
    go self.write()
    go self.read()
}

// ---- write begin ----
func (self *ClientSession) write() {
    // clean func
    defer func() {   
        self.Close()     
    }()

    for {
        select {
        case pWrite := <-self.chSend:
            // pWrite中的错误信息不会发到客户端
            // close agent if low-level Conn broken
            if _, err := self.conn.Write(pWrite.data); err != nil {                
                logger.Log.Errorf("Failed to write in conn: %s", err.Error())
                return
            }            
        case <- self.chanClose:
            return
        }
    }
}

func (self *ClientSession) getMessageFromPendingMessage(pm pendingMessage) (*message.Message, error) {
    payload, ok :=  pm.payload.([]byte)
    if !ok {
        // TODO: 错误
        return nil, errors.New("bad payload")
    }

    // construct message and encode
    m := &message.Message{
        Type:  pm.typ,
        Data:  payload,
        Route: pm.route,
        ID:    pm.mid,
        Err:   pm.err,
    }

    return m, nil
}

func (self *ClientSession) send(pendingMsg pendingMessage) (err error) {
    defer func() {
        if e := recover(); e != nil {
            err = errors.New("send error")
        }
    }()    

    m, err := self.getMessageFromPendingMessage(pendingMsg)
    if err != nil {
        return err
    }

    // packet encode
    p, err := self.packetEncodeMessage(m)
    if err != nil {
        return err
    }

    pWrite := &pendingWrite{        
        data: p,
    }

    if pendingMsg.err {
        pWrite.err = errors.New("getMessageFromPendingMessage error")
    }

    // chSend is never closed so we need this to don't block if agent is already closed
    select {
    case self.chSend <- pWrite:
    //case <-self.chDie:
    }
    return
}

func (self *ClientSession) packetEncodeMessage(m *message.Message) ([]byte, error) {
    em, err := self.cfg.MessageEncoder.Encode(m)
    if err != nil {
        return nil, err
    }

    // packet encode
    p, err := self.cfg.Encoder.Encode(packet.Data, em)
    if err != nil {
        return nil, err
    }
    return p, nil
}

// Push implementation for NetworkEntity interface
func (self *ClientSession) Push(route string, v interface{}) error {
    if self.GetStatus() == StatusClosed {
        return errors.New("closed")
    }

    // switch d := v.(type) {
    // case []byte:
    //     logger.Log.Debugf("Type=Push, ID=%d, UID=%s, Route=%s, Data=%dbytes",
    //         self.Session.ID(), a.Session.UID(), route, len(d))
    // default:
    //     logger.Log.Debugf("Type=Push, ID=%d, UID=%s, Route=%s, Data=%+v",
    //         a.Session.ID(), a.Session.UID(), route, v)
    // }
    return self.send(pendingMessage{typ: message.Push, route: route, payload: v})
}

// SendHandshakeResponse sends a handshake response
func (self *ClientSession) SendHandshakeResponse() error {
    _, err := self.conn.Write(hrd)
    return err
}


// ResponseMID implementation for NetworkEntity interface
// Respond message to session
func (self *ClientSession) ResponseMID(mid uint, v interface{}, e error) error {
    err := false
    if e != nil {
        err = true
    }
    if self.GetStatus() == StatusClosed {
        return errors.New("closed")
    }

    if mid <= 0 {
        return constants.ErrSessionOnNotify
    }

    // switch d := v.(type) {
    // case []byte:        
    //     logger.Log.Debugf("Type=Response, ID=%d, UID=%s, MID=%d, Data=%dbytes",
    //          0, 0, mid, len(d))
    // default:
    //     logger.Log.Infof("Type=Response, ID=%d, UID=%s, MID=%d, Data=%+v",
    //         0, 0, mid, v)
    // }

    return self.send(pendingMessage{typ: message.Response, mid: mid, payload: v, err: err})
}

// ---- write over ----
// ---- read begin ----
func (self *ClientSession) read() {  
    conn := self.conn
    for {
        msg, err := conn.GetNextMessage()

        if err != nil {
            if err != constants.ErrConnectionClosed {
                logger.Log.Errorf("Error reading next available message: %s", err.Error())
            }
            self.Close()
            return
        }

        packets, err := self.cfg.Decoder.Decode(msg)
        if err != nil {
            logger.Log.Errorf("Failed to decode message: %s", err.Error())
            return
        }

        if len(packets) < 1 {
            logger.Log.Warnf("Read no packets, data: %v", msg)
            continue
        }

        // process all packet
        for i := range packets {
            if err := self.processPacket(packets[i]); err != nil {
                logger.Log.Errorf("Failed to process packet: %s", err.Error())
                return
            }
        }
    }
}

func (self *ClientSession) processPacket(p *packet.Packet) error {
    switch p.Type {
    case packet.Handshake:
        logger.Log.Debug("Received handshake packet")
        if err := self.SendHandshakeResponse(); err != nil {
            logger.Log.Errorf("Error sending handshake response: %s", err.Error())
            return err
        }
        //logger.Log.Debugf("Session handshake Id=%d, Remote=%s", a.GetSession().ID(), a.RemoteAddr())

        // Parse the json sent with the handshake by the client
        handshakeData := &HandshakeData{}
        err := json.Unmarshal(p.Data, handshakeData)
        if err != nil {
            self.SetStatus(StatusClosed)
            return fmt.Errorf("Invalid handshake data. Id=%d", 0)
        }

        //a.GetSession().SetHandshakeData(handshakeData)
        self.SetStatus(StatusHandshake)
        //err = a.GetSession().Set(constants.IPVersionKey, a.IPVersion())
        // if err != nil {
        //     logger.Log.Warnf("failed to save ip version on session: %q\n", err)
        // }

        logger.Log.Debug("Successfully saved handshake data")

    case packet.HandshakeAck:
        self.lastHeartBeat = common.NowMs()
        self.SetStatus(StatusWorking)
        //logger.Log.Debugf("Receive handshake ACK Id=%d, Remote=%s", a.GetSession().ID(), a.RemoteAddr())

    case packet.Data:
        if self.GetStatus() < StatusWorking {
            // return fmt.Errorf("receive data on socket which is not yet ACK, session will be closed immediately, remote=%s",
            //     a.RemoteAddr().String())
            return nil
        }

        msg, err := message.Decode(p.Data)
        if err != nil {
            return err
        }
        self.processMessage(msg)
        //logger.Log.Debugf("%v", msg)

    // 心跳处理
    // 主要是断开不活跃的连接
    case packet.Heartbeat:
        //logger.Log.Debugf("got heartbeat")
        self.lastHeartBeat = common.NowMs()
    }

    //a.SetLastAt()
    return nil
}
// ---- read over ----
// ---- heartbeat ----
func (self *ClientSession) heartbeat() {
    // timer
    t := time.NewTicker(DefaultHeartbeatTime)
    defer func() {
        t.Stop()
    }()

    for {
        select {    
        case <- t.C:
            self.sendHeartBeat()            
            self.checkHeartBeatTimeout()
        case <- self.chanClose:
            return
        }
    }
}

func (self *ClientSession) checkHeartBeatTimeout() {
    if self.GetStatus() != StatusWorking {
        return
    }
    if common.NowMs() < self.lastHeartBeat + 2 * DefaultHeartbeatTimeSeconds*1000 {        
        return
    }

    logger.Log.Debugf("heartbeat timeout and close")
    // close self
    self.Close()
}

func (self *ClientSession) sendHeartBeat() {
    if self.GetStatus() != StatusWorking {
        return
    }
    pWrite := &pendingWrite{        
        data: hbd,
    }    

    // chSend is never closed so we need this to don't block if agent is already closed
    select {
    case self.chSend <- pWrite:
    //case <-self.chDie:
    }
}
// ---- heartbeat over ----

func (self *ClientSession) processMessage(msg *message.Message) {
    self.GetImpl().ProcessMessage(self, msg)    
}


// 生成一个hdd packet
func hbdEncode(heartbeatTimeout time.Duration, packetEncoder codec.PacketEncoder, dataCompression bool, serializerName string) {
    hData := map[string]interface{}{
        "code": 200,
        "sys": map[string]interface{}{
            "heartbeat":  heartbeatTimeout.Seconds(),
            "dict":       message.GetDictionary(),
            "serializer": serializerName,
        },
    }
    data, err := json.Marshal(hData)
    if err != nil {
        panic(err)
    }

    if dataCompression {
        logger.Log.Debugf("hbdEncode compression")
        compressedData, err := compression.DeflateData(data)
        if err != nil {
            panic(err)
        }

        if len(compressedData) < len(data) {
            data = compressedData
        }
    }

    hrd, err = packetEncoder.Encode(packet.Handshake, data)
    if err != nil {
        panic(err)
    }

    hbd, err = packetEncoder.Encode(packet.Heartbeat, nil)
    if err != nil {
        panic(err)
    }
}
