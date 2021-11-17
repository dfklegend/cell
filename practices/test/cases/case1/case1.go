package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/dfklegend/cell/practices/test/network/tcp/server"
	"github.com/dfklegend/cell/practices/test/routines"
)

// 测试ReadFull是否有空转
// 结论是:没有，Read会阻塞住
func main() {
	fmt.Println("-----")
	go server.ListenAndServe()
	routines.RunCounter(32)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}
