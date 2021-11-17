package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	str := "Hello, World!"
	fmt.Println(str)
	

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}
