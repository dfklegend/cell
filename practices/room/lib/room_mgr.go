package lib

import (
	"fmt"
	"runtime"
	"time"
)

type RoomMgr struct {
	rooms map[int]*Room
	tick  *time.Ticker

	nextId     int
	maxRoomNum int
}

func (self *RoomMgr) AddRoom(room *Room) {
	self.rooms[room.GetId()] = room
}

// 如果每次创建一个房间
func (self *RoomMgr) Update() {
	size := len(self.rooms)
	if size < self.maxRoomNum {
		self.NewRoom()
	}

	self.UpdateRooms()
}

func (self *RoomMgr) allocId() int {
	id := self.nextId
	self.nextId += 1
	return id
}

func (self *RoomMgr) NewRoom() {
	id := self.allocId()
	room := new(Room)
	room.Init(id)
	self.AddRoom(room)
	room.Start()

	runtime.SetFinalizer(room, func(r *Room) {
		fmt.Printf("finalizer room:%d\n", r.GetId())
	})

	fmt.Printf("new room:%d\n", id)
}

func (self *RoomMgr) UpdateRooms() {
	for id, room := range self.rooms {
		if room.IsOver() {
			delete(self.rooms, id)
			fmt.Printf("remove room:%d\n", id)
		}
	}
}

func (self *RoomMgr) Start(num int) {
	self.tick = time.NewTicker(100 * time.Millisecond)

	self.rooms = make(map[int]*Room)
	self.maxRoomNum = num
	go self.loop()
}

func (self *RoomMgr) loop() {
	for {
		select {
		case <-self.tick.C:
			self.Update()
		}
	}
}
