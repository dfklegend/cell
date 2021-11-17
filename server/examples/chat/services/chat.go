package services

import (
	"fmt"
	"log"
	"sort"
	"math/rand"

	api "github.com/dfklegend/cell/rpc/apientry"
	"github.com/dfklegend/cell/utils/logger"
	ucommon "github.com/dfklegend/cell/utils/common"
	"github.com/dfklegend/cell/utils/runservice"

	"github.com/dfklegend/cell/server/client/session"
	"github.com/dfklegend/cell/server/interfaces"
	"github.com/dfklegend/cell/server/common"
	"github.com/dfklegend/cell/server/services/channel"
	"github.com/dfklegend/cell/server/examples/chat/protos"
)

const (
	MaxPlayerInRoom = 3
)

// handler
type Chat struct {
	api.APIEntry
}

func (self *Chat) SendChat(d api.IHandlerSession, msg *protos.ChatMsg, cbFunc api.HandlerCBFunc) error {
	logger.Log.Debugf("enter Chat.SendChat")

	bs := d.(*session.BackSession)
	
	cs := GetChatService()
	cs.OnChat(bs.GetID(), msg.Content)
	// r := cs.GetPlayerRoom(bs.GetID())
	// if r != nil {
	// 	r.GetChannel().PushMessage("onMessage", msg)
	// }

	// 请求进入
	api.CheckInvokeCBFunc(cbFunc, nil,
		&protos.NormalAck{
			Code: 0,
			Result: "hello ack",
			})
	return nil
}

// remote
type ChatRemote struct {
	api.APIEntry
}

func (self *ChatRemote) Entry(msg *protos.RoomEntry, cbFunc api.HandlerCBFunc) error {
	logger.Log.Debugf("enter ChatRemote.Entry:%+v", msg)		

	// test
	app := interfaces.GetApp()	
	frontId := msg.ServerId
	netId := msg.NetId
	app.GetSysService().PushMessageById(frontId, netId, "onTest", "{\"a\":1}", nil)
	//app.GetSysService().PushMessageById(frontId, netId, "onTest", msg, nil)
	
	// cs := channel.GetChannelService()	
	// // 先给自己发members

	// c := cs.AddToChannel("chat", frontId, netId)
	// newMsg := &protos.OnNewUser{
	// 	Name: fmt.Sprintf("%v.%v", frontId, netId),
	// }
	// c.PushMessage("onNewUser", newMsg)
	chat := GetChatService()
	chat.PlayerEnter(msg.UId, msg.Name, msg.ServerId, msg.NetId)

	// 请求进入
	api.CheckInvokeCBFunc(cbFunc, nil,
		&protos.NormalAck{
			Code: 0,
			Result: "n ack",
			})
	return nil
}

func (self *ChatRemote) Leave(msg *protos.RoomLeave, cbFunc api.HandlerCBFunc) error {
	
	frontId := msg.ServerId
	netId := msg.NetId

	log.Printf("leave %v,%v\n", frontId, netId)

	// cs := channel.GetChannelService()
	// cs.LeaveFromChannel("chat", frontId, netId)

	// c := cs.GetChannel("chat")	
	// if c != nil {
	// 	leaveMsg := &protos.OnUserLeave {
	// 		Name: fmt.Sprintf("%v.%v", frontId, netId),
	// 	}
	// 	c.PushMessage("onUserLeave", leaveMsg)	
	// }
	chat := GetChatService()
	chat.PlayerLeave(msg.UId)
	

	api.CheckInvokeCBFunc(cbFunc, nil,
		&protos.NormalAck{
			Code: 0,
			Result: "n ack",
	})
	return nil
}

// ----
var chatService *ChatService = NewChatService()

// 聊天具体的service
// 按登录 2人分一个房间
// 聊天房间内同步
// 输入roll,掷骰子
// 离开后，其他人可以加入
type ChatService struct {
	// 不同房间
	runService *runservice.RunService	
	idService *ucommon.SerialIdService

	rooms map[uint32] *ChatRoom
	players map[string] *ChatPlayer
}

func NewChatService() *ChatService {
	service := &ChatService {
		runService: runservice.NewRunService("chatservice"),
		idService: ucommon.NewSerialIdService(),
		rooms: make(map[uint32] *ChatRoom),
		players: make(map[string] *ChatPlayer),
	}
	return service
}

func GetChatService() *ChatService {
	return chatService
}

func (self* ChatService) Start() {
	self.runService.Start()
}

func (self *ChatService) findPlayer(uid string) *ChatPlayer {
	p, _ := self.players[uid]
	return p
}

func (self *ChatService) findRoom(roomId uint32) *ChatRoom {
	r, _ := self.rooms[roomId]
	return r
}

// TODO: 优先房间id比较小得
func (self *ChatService) findEmptyRoom() *ChatRoom {
	keys := make([]uint32, 0, len(self.rooms))

	for k, _ := range(self.rooms) {
		keys = append(keys, k)
	}

	//sort.Ints(keys)
	sort.SliceStable(keys, func(a, b int) bool {
		return keys[a] < keys[b]
	})

	for _, k := range(keys) {
		v := self.findRoom(k)
		if v.GetPlayerNum() < MaxPlayerInRoom {
			return v
		}
	}
	return nil
}

func (self *ChatService) GetPlayerRoom(uid string) *ChatRoom {
	p := self.findPlayer(uid)
	if p == nil {
		return nil
	}
	room := self.findRoom(p.RoomId)
	if room == nil {
		return nil
	}
	return room
}

func (self* ChatService) CreateRoom() *ChatRoom {
	id := self.idService.AllocId()
	self.rooms[id] = NewRoom(id)
	return self.rooms[id]
}

func (self* ChatService) PlayerEnter(uid string, name string,
	frontId string, netId common.NetIdType) {
	p := NewPlayer()
	p.UId = uid
	p.Name = name
	p.FrontId = frontId
	p.NetId = netId

	self.players[uid] = p	

	r := self.findEmptyRoom()
	if r == nil {
		r = self.CreateRoom()
	}

	// 推送房间members
	self.pushMembers(p, r)

	r.AddPlayer(p)
	// 广播进入房间	
	newMsg := &protos.OnNewUser{
		Name: p.GetDetailName(),
	}
	r.GetChannel().PushMessage("onNewUser", newMsg)

	// 显示系统信息
	self.showWelcome(p)
}

func (self *ChatService) pushMembers(p *ChatPlayer, r *ChatRoom) {
	cs := channel.GetChannelService()	

	members := r.GetMembers()
	if len(members) == 0 {
		return
	}
	
	cs.PushMessageById(p.FrontId, p.NetId, "onMembers",
		&protos.OnMembers{
			Members: members,
	}, nil)
}

func (self *ChatService) showWelcome(p *ChatPlayer) {
	cs := channel.GetChannelService()
	app := interfaces.GetApp()
	str := fmt.Sprintf("welcome to %v room:%v", app.GetServerId(), p.RoomId)
	logger.Log.Debugf(str)
	cs.PushMessageById(p.FrontId, p.NetId, "onMessage",
		&protos.ChatMsg{
			Name: "system",
			Content: str,
		}, nil)
}


func (self* ChatService) PlayerLeave(uid string) {
	// find player
	p := self.findPlayer(uid)
	if p == nil {
		return
	}
	room := self.findRoom(p.RoomId)
	if room == nil {
		return
	}

	room.RemovePlayer(uid)
	leaveMsg := &protos.OnUserLeave {
		Name: p.GetDetailName(),
	}
	room.GetChannel().PushMessage("onUserLeave", leaveMsg)	

	delete(self.players, uid)
}

func (self *ChatService) OnChat(uid, content string) {
	p := self.findPlayer(uid)
	if p == nil {
		return
	}
	room := self.findRoom(p.RoomId)
	if room == nil {
		return
	}	
	room.OnChat(p, content)
}

// ----

type ChatPlayer struct {
	UId string
	Name string
	FrontId string
	NetId common.NetIdType
	RoomId uint32
}

func NewPlayer() *ChatPlayer {
	return &ChatPlayer{}
}

func (self* ChatPlayer) GetDetailName() string {
	return fmt.Sprintf("%v %v.%v", self.Name, self.FrontId, self.NetId)
}

type ChatRoom struct {
	id uint32
	players map[string]*ChatPlayer
	c *channel.Channel
}

func NewRoom(id uint32) *ChatRoom {
	cs := channel.GetChannelService()
	return &ChatRoom{
		id: id,
		players: make(map[string]*ChatPlayer),
		c: cs.AddChannel(fmt.Sprintf("room%v", id)),
	}
}

func (self *ChatRoom) Destroy() {
	cs := channel.GetChannelService()
	cs.DeleteChannel(self.c.GetName())
}

func (self *ChatRoom) AddPlayer(p *ChatPlayer) {	
	self.c.Add(p.FrontId, p.NetId)
	self.players[p.UId] = p
	p.RoomId = self.id
}

func (self *ChatRoom) RemovePlayer(uid string) {
	p, _ := self.players[uid]
	if p == nil {
		return
	}
	self.c.Leave(p.FrontId, p.NetId)
	delete(self.players, uid)
	p.RoomId = 0
}

func (self *ChatRoom) GetChannel() *channel.Channel {
	return self.c
}

func (self *ChatRoom) GetPlayerNum() int {
	return len(self.players)
}

func (self *ChatRoom) GetMembers() []string {
	members := make([]string, 0)
	for _, v := range(self.players) {
		members = append(members, v.GetDetailName())
	}
	return members
}

func (self *ChatRoom) OnChat(p *ChatPlayer, content string) {
	logger.Log.Debugf("OnChat:", content)
	if content == "/roll" {		
		// do roll
		str := fmt.Sprintf("掷出了%v点(1-100)",
			rand.Intn(100))
		self.PushMessage(p, str)
		return
	}
	self.PushMessage(p, content)
}

func (self *ChatRoom) PushMessage(p *ChatPlayer, content string) {	
	self.c.PushMessage("onMessage", &protos.ChatMsg{
		Name: p.Name,
		Content: content,
	})
}
