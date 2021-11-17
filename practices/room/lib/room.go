package lib

import (
	"fmt"
	"math/rand"
	"time"
)

type Room struct {
	id   int
	over bool
	tick *time.Ticker

	runTicks int
}

func (self *Room) Update() {
	// 每秒输出一次
	// 每秒1%几率over
	self.runTicks++
	if self.runTicks%10 == 0 {
		//fmt.Printf("id:%d ticks:%d\n", self.id, self.runTicks)

		if rand.Float32() <= 0.01 {
			self.SetOver()
		}
	}
}

func (self *Room) GetId() int {
	return self.id
}

func (self *Room) SetOver() {
	fmt.Printf("id:%d setOver\n", self.id)
	self.over = true
}

func (self *Room) IsOver() bool {
	return self.over
}

func (self *Room) Init(id int) {
	self.id = id
}

func (self *Room) Start() {
	self.tick = time.NewTicker(20 * time.Millisecond)
	go self.loop()
}

func (self *Room) loop() {
	for !self.IsOver() {
		select {
		case <-self.tick.C:
			self.Update()
		}
	}
}
