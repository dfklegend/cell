package main

import (
	"fmt"
	"os"
	"os/signal"

	"dfk.com/practices/room/lib"
)

func main() {
	str := "Hello, World!"
	fmt.Println(str)

	mgr := new(lib.RoomMgr)
	mgr.Start(100)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}
