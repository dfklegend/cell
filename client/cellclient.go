

package client

import (	
	"time"
	"reflect"

	"github.com/dfklegend/cell/utils/logger"
	"github.com/dfklegend/cell/utils/sche"
	"github.com/dfklegend/cell/utils/runservice"	
	nclient "github.com/dfklegend/cell/net/client"
)

const (
	STATE_INIT = iota
	STATE_CONNECTING
	STATE_CONNECTED
	STATE_BROKEN
	STATE_PENDINGRETRY 		// 等待5s再重连
)

type HandleFunc func()

// add auto retry
type CellClient struct {
	TheClient	*nclient.Client
	runService *runservice.StandardRunService
	state int
	firstStart bool
	autoRetry bool

	connectingTime int
	retryWait int
	tarAddress string

	cbBreak HandleFunc 
	cbConnected HandleFunc
}

func NewCellClient(name string) *CellClient {
	return &CellClient {
		TheClient: nclient.New(),
		runService: runservice.NewStandardRunService(name),
		firstStart: true,
		autoRetry: true,
		state: STATE_INIT,
	}
}

func (self *CellClient) GetRunService() *runservice.StandardRunService {
	return self.runService
}

func (self *CellClient) setState(state int) {
	self.state = state
}

func (self *CellClient) getState() int {
	return self.state
}

func (self *CellClient) IsReady() bool {
	return self.getState() == STATE_CONNECTED
}

func (self *CellClient) SetCBBreak(cb HandleFunc) {
	self.cbBreak = cb
}

func (self *CellClient) SetCBConnected(cb HandleFunc) {
	self.cbConnected = cb
}

func (self *CellClient) Start(address string) {	
	if self.firstStart {
		self.runService.Start()
		self.runService.GetEventCenter().SetLocalUseChan(true)	
		self.addUpdate()

		self.firstStart = false
	}
	self.setState(STATE_INIT)
	self.autoRetry = true

	self.tarAddress = address
	self.Connect(address)
}

func (self *CellClient) Connect(address string) {
	old := self.TheClient
	if old != nil {
		old.Disconnect()
	}
	self.TheClient = nclient.New()
	self.setState(STATE_CONNECTING)
	self.TheClient.ConnectTo(address)

	// TODO: 移除掉老的MsgChannel selector
	self.addMsgProcess()
}

func (self *CellClient) WaitReady() {
	for !self.TheClient.Ready {
		time.Sleep(100*time.Millisecond)
	}
}

func (self *CellClient) Stop() {	
	self.autoRetry = false
	self.TheClient.Disconnect()	
	self.runService.Stop()
}

func (self *CellClient) GetClient() *nclient.Client {
	return self.TheClient
}

func (self *CellClient) addMsgProcess() {
	selector := self.runService.GetSelector()
	selector.AddSelector(sche.NewFuncSelector(reflect.ValueOf(self.TheClient.MsgChannel()),
		func(v reflect.Value, recvOk bool) {
			if !recvOk {
				return
			}

			msg := v.Interface().(*nclient.ClientMsg)
			
			if msg.Cb != nil {
				msg.Cb(false, msg.Msg)
			} else {
				// push msg
				logger.Log.Debugf("got push:%v", msg.Msg)
				Msg := msg.Msg				
				self.runService.GetEventCenter().Publish(Msg.Route, Msg.Data)
			}
	}))
}

func (self *CellClient) addUpdate() {
	self.runService.GetTimerMgr().AddTimer(1*time.Second, func(args ...interface{}) {
		self.onUpdate()
	})
}

func (self *CellClient) onUpdate() {
	switch(self.getState()) {
	case STATE_CONNECTING:
		self.onConnecting()
	case STATE_CONNECTED:
		if !self.TheClient.Connected {
			self.setState(STATE_BROKEN)
			self.onBreak()
		}
	case STATE_PENDINGRETRY:
		self.onPendingRetry()
	}
}

func (self *CellClient) onConnecting() {
	if self.TheClient.Ready {
		self.setState(STATE_CONNECTED)
		if self.cbConnected != nil {
			self.cbConnected()
		}
		return
	}
	self.connectingTime ++
	if self.connectingTime > 5 {
		self.onBreak()
	}
}

func (self *CellClient) onBreak() {
	if self.cbBreak != nil {
		self.cbBreak()
	}
	// 看要不要重连
	if self.autoRetry {
		self.beginRetry()
	}
}

func (self *CellClient) beginRetry() {
	self.setState(STATE_PENDINGRETRY)
	self.retryWait = 5	
}

func (self *CellClient) onPendingRetry() {
	self.retryWait --
	logger.Log.Debugf("reconnect wait:%v", self.retryWait)
	if self.retryWait <= 0 {
		logger.Log.Debugf("reconnect")
		self.Connect(self.tarAddress)
	}
}