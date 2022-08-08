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

package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dfklegend/cell/net/common/conn/codec"
	"github.com/dfklegend/cell/net/common/conn/message"
	"github.com/dfklegend/cell/net/common/conn/packet"
	"github.com/dfklegend/cell/utils/compression"
	"github.com/dfklegend/cell/utils/logger"
)

const (
	ReqeustTimeout = 15 * time.Second
)

var (
	handshakeBuffer = `
{
	"sys": {
		"platform": "mac",
		"libVersion": "0.3.5-release",
		"clientBuildNumber":"20",
		"clientVersion":"2.1"
	},
	"user": {
		"age": 30
	}
}
`
)

// HandshakeSys struct
type HandshakeSys struct {
	Dict       map[string]uint16 `json:"dict"`
	Heartbeat  int               `json:"heartbeat"`
	Serializer string            `json:"serializer"`
}

// HandshakeData struct
type HandshakeData struct {
	Code int          `json:"code"`
	Sys  HandshakeSys `json:"sys"`
}

type ClientMsg struct {
	Msg *message.Message
	Cb  CBTask
}

// 加入一个callback
type pendingRequest struct {
	msg    *message.Message
	sentAt time.Time
	cb     CBTask
}

type CBTask func(error bool, msg *message.Message)

// Client
// 重连在上层处理
// 携程
//     readServerMessages
//     		读取消息
//     		组织成packet
//     handlePackets
//     		处理packet
//     sendHeartbeats
//     		定时发送heartbeat
//     pendingRequestsReaper
//     		移除超时未返回的请求
type Client struct {
	conn      net.Conn
	Connected bool
	// 握手完毕
	Ready           bool
	packetEncoder   codec.PacketEncoder
	packetDecoder   codec.PacketDecoder
	packetChan      chan *packet.Packet
	IncomingMsgChan chan *ClientMsg
	pendingChan     chan bool
	pendingRequests map[uint]*pendingRequest
	pendingReqMutex sync.Mutex
	requestTimeout  time.Duration
	closeChan       chan struct{}
	nextID          uint32
	messageEncoder  message.Encoder
}

// MsgChannel return the incoming message channel
func (c *Client) MsgChannel() chan *ClientMsg {
	return c.IncomingMsgChan
}

// ConnectedStatus return the connection status
func (c *Client) ConnectedStatus() bool {
	return c.Connected
}

// New returns a new client
func New(requestTimeout ...time.Duration) *Client {

	reqTimeout := ReqeustTimeout
	if len(requestTimeout) > 0 {
		reqTimeout = requestTimeout[0]
	}

	return &Client{
		Connected:       false,
		Ready:           false,
		packetEncoder:   codec.NewPomeloPacketEncoder(),
		packetDecoder:   codec.NewPomeloPacketDecoder(),
		packetChan:      make(chan *packet.Packet, 10),
		pendingRequests: make(map[uint]*pendingRequest),
		requestTimeout:  reqTimeout,
		// 30 here is the limit of inflight messages
		// TODO this should probably be configurable
		pendingChan:    make(chan bool, 30),
		messageEncoder: message.NewMessagesEncoder(false),
	}
}

func (c *Client) onError() {
}

// Disconnect disconnects the client
func (c *Client) Disconnect() {
	if c.Connected {
		c.Connected = false
		c.Ready = false
		close(c.closeChan)
		c.conn.Close()
	}
}

// ConnectToTLS connects to the server at addr using TLS, for now the only supported protocol is tcp
// this methods blocks as it also handles the messages from the server
func (c *Client) ConnectToTLS(addr string, skipVerify bool) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: skipVerify,
	})
	if err != nil {
		return err
	}
	return c.start(conn)
}

// ConnectTo connects to the server at addr, for now the only supported protocol is tcp
// this methods blocks as it also handles the messages from the server
func (c *Client) ConnectTo(addr string) error {
	fmt.Printf("begin dial:%v\n", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("dial error:%v\n", err)
		return err
	}
	return c.start(conn)
}

func (c *Client) start(conn net.Conn) error {
	c.conn = conn
	c.IncomingMsgChan = make(chan *ClientMsg, 10)

	fmt.Printf("dial over\n")

	// 等一会
	//fmt.Println("wait 5 seconds")
	//time.Sleep(5 * time.Second)
	if err := c.beginHandshake(); err != nil {
		return err
	}

	c.closeChan = make(chan struct{})

	return nil
}

func (c *Client) sendHandshakeRequest() error {
	p, err := c.packetEncoder.Encode(packet.Handshake, []byte(handshakeBuffer))
	if err != nil {
		return err
	}
	_, err = c.conn.Write(p)
	return err
}

func (c *Client) beginRecvMessages() error {
	c.Connected = true

	go c.readServerMessages()
	go c.handlePackets()

	return nil
}

// 等到一个握手返回，整个连接就建立成功
func (c *Client) handleHandleShake(handshakePacket *packet.Packet) error {
	fmt.Println("got handleHandleShake")
	if handshakePacket.Type != packet.Handshake {
		return fmt.Errorf("got first packet from server that is not a handshake, aborting")
	}

	var err error
	handshake := &HandshakeData{}
	if compression.IsCompressed(handshakePacket.Data) {
		handshakePacket.Data, err = compression.InflateData(handshakePacket.Data)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(handshakePacket.Data, handshake)
	if err != nil {
		return err
	}

	logger.Log.Debug("got handshake from sv, data: %v", handshake)

	if handshake.Sys.Dict != nil {
		message.SetDictionary(handshake.Sys.Dict)
	}
	p, err := c.packetEncoder.Encode(packet.HandshakeAck, []byte{})
	if err != nil {
		return err
	}
	_, err = c.conn.Write(p)
	if err != nil {
		return err
	}

	logger.Log.Debug("connection ready")
	// 握手成功
	c.Ready = true
	go c.sendHeartbeats(handshake.Sys.Heartbeat)
	go c.pendingRequestsReaper()
	return nil
}

// pendingRequestsReaper delete timedout requests
func (c *Client) pendingRequestsReaper() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			toDelete := make([]*pendingRequest, 0)
			c.pendingReqMutex.Lock()
			for _, v := range c.pendingRequests {
				if time.Now().Sub(v.sentAt) > c.requestTimeout {
					toDelete = append(toDelete, v)
				}
			}
			for _, pendingReq := range toDelete {

				logger.Log.Infof("request timeout: %v", pendingReq.msg.ID)
				err := errors.New("request timeout")
				errMarshalled, _ := json.Marshal(err)

				// send a timeout to incoming msg chan
				m := &message.Message{
					Type:  message.Response,
					ID:    pendingReq.msg.ID,
					Route: pendingReq.msg.Route,
					Data:  errMarshalled,
					Err:   true,
				}
				nm := &ClientMsg{
					Msg: m,
					Cb:  nil,
				}

				var cbTask CBTask
				pendingID := pendingReq.msg.ID
				if req, ok := c.pendingRequests[pendingID]; ok {
					cbTask = req.cb
					delete(c.pendingRequests, pendingID)
					<-c.pendingChan
				}
				nm.Cb = cbTask
				c.IncomingMsgChan <- nm
			}
			c.pendingReqMutex.Unlock()
		case <-c.closeChan:
			return
		}
	}
}

// 处理包
func (c *Client) handlePackets() {
	for {
		select {
		case p := <-c.packetChan:
			switch p.Type {
			case packet.Handshake:
				c.handleHandleShake(p)
			case packet.Data:
				//handle data
				//logger.Log.Debug("got data: %s", string(p.Data))
				var cbTask CBTask
				m, err := message.Decode(p.Data)
				if err != nil {
					logger.Log.Errorf("error decoding msg from sv: %s", string(m.Data))
				}
				if m.Type == message.Response {
					c.pendingReqMutex.Lock()
					if req, ok := c.pendingRequests[m.ID]; ok {
						cbTask = req.cb
						delete(c.pendingRequests, m.ID)
						<-c.pendingChan
					} else {
						c.pendingReqMutex.Unlock()
						continue // do not process msg for already timedout request
					}
					c.pendingReqMutex.Unlock()
				}
				//fmt.Println("cbTask")
				//fmt.Println(cbTask)
				//if cbTask != nil {
				//	cbTask(false, m)
				//}
				nm := &ClientMsg{
					Msg: m,
					Cb:  cbTask,
				}
				c.IncomingMsgChan <- nm
			case packet.Kick:
				logger.Log.Warn("got kick packet from the server! disconnecting...")
				c.Disconnect()
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) readPackets(buf *bytes.Buffer) ([]*packet.Packet, error) {
	// listen for sv messages
	data := make([]byte, 1024)
	n := len(data)
	var err error

	for n == len(data) {
		n, err = c.conn.Read(data)
		//fmt.Printf("conn.read ret:%d %v", n, err)
		if err != nil {
			return nil, err
		}
		buf.Write(data[:n])
	}
	packets, err := c.packetDecoder.Decode(buf.Bytes())
	if err != nil {
		logger.Log.Errorf("error decoding packet from server: %s", err.Error())
	}
	totalProcessed := 0
	for _, p := range packets {
		totalProcessed += codec.HeadLength + p.Length
	}
	buf.Next(totalProcessed)

	return packets, nil
}

// 读取包
func (c *Client) readServerMessages() {
	buf := bytes.NewBuffer(nil)
	defer c.Disconnect()
	for c.Connected {
		packets, err := c.readPackets(buf)
		if err != nil && c.Connected {
			logger.Log.Error(err)
			break
		}

		for _, p := range packets {
			c.packetChan <- p
		}
	}
}

func (c *Client) sendHeartbeats(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer func() {
		t.Stop()
		c.Disconnect()
	}()
	for {
		select {
		case <-t.C:
			p, _ := c.packetEncoder.Encode(packet.Heartbeat, []byte{})
			_, err := c.conn.Write(p)
			if err != nil {
				logger.Log.Errorf("error sending heartbeat to server: %s", err.Error())
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) beginHandshake() error {
	if err := c.sendHandshakeRequest(); err != nil {
		return err
	}

	if err := c.beginRecvMessages(); err != nil {
		return err
	}
	return nil
}

// 等待连接完毕
func (c *Client) WaitReady() {
	for !c.Ready {
		time.Sleep(100 * time.Millisecond)
	}
}

// SendRequest sends a request to the server
func (c *Client) SendRequest(route string, data []byte, cb CBTask) (uint, error) {
	return c.sendMsg(message.Request, route, data, cb)
}

// SendNotify sends a notify to the server
func (c *Client) SendNotify(route string, data []byte) error {
	_, err := c.sendMsg(message.Notify, route, data, nil)
	return err
}

func (c *Client) buildPacket(msg message.Message) ([]byte, error) {
	encMsg, err := c.messageEncoder.Encode(&msg)
	if err != nil {
		return nil, err
	}
	p, err := c.packetEncoder.Encode(packet.Data, encMsg)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// sendMsg sends the request to the server
func (c *Client) sendMsg(msgType message.Type, route string, data []byte, cb CBTask) (uint, error) {
	// TODO mount msg and encode
	if !c.Ready {
		logger.Log.Errorf("sendMsg error, client not Ready: %v", route)
		return 0, nil
	}

	m := message.Message{
		Type:  msgType,
		ID:    uint(atomic.AddUint32(&c.nextID, 1)),
		Route: route,
		Data:  data,
		Err:   false,
	}
	p, err := c.buildPacket(m)
	if msgType == message.Request {
		c.pendingChan <- true
		c.pendingReqMutex.Lock()
		if _, ok := c.pendingRequests[m.ID]; !ok {
			c.pendingRequests[m.ID] = &pendingRequest{
				msg:    &m,
				sentAt: time.Now(),
				cb:     cb,
			}
		}
		c.pendingReqMutex.Unlock()
	}

	if err != nil {
		return m.ID, err
	}

	//fmt.Printf("conn:%v\n", c.conn)
	_, err = c.conn.Write(p)
	if err != nil {
		c.onError()
	}
	return m.ID, err
}
