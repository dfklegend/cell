package main

import (
	"log"	
)

type Base struct {	
}

func (self *Base) Start() {
	log.Println("Base.Start")
}

type Child struct {
	Base
}

func (self *Child) Start() {
	self.Base.Start()
	log.Println("Child.Start")	
}

func main() {
	a := Child{}
	a.Start()
}