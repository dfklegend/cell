package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	scheP "github.com/dfklegend/cell/practices/test/libs/sche"
	"github.com/dfklegend/cell/practices/test/network/tcp/server"
)

/*
1. routine1
	延迟2秒，模拟载入玩家请求
2. mgr routine1
	实际处理载入玩家
3. db routine1
	响应载入玩家数据请求
*/
func main() {
	fmt.Println("-----")
	go server.ListenAndServe()

	go TheMgr.Handle()
	go TheDB.Handle()

	// go func() {
	// 	time.Sleep(2 * time.Second)
	// 	var sche = scheP.TheScheMgr.GetSche("mgr")
	// 	sche.Post(func() (ret interface{}, err error) {
	// 		TheMgr.LoadPlayer(0)
	// 		return 0, nil
	// 	})
	// }()
	// go func() {
	// 	time.Sleep(3 * time.Second)
	// 	var sche = scheP.TheScheMgr.GetSche("mgr")
	// 	sche.Post(func() (ret interface{}, err error) {
	// 		TheMgr.LoadPlayer(1)
	// 		return 0, nil
	// 	})
	// }()

	for i := 0; i < 10; i++ {
		uid := i
		go func() {
			time.Sleep(3 * time.Second)
			var sche = scheP.TheScheMgr.GetSche("mgr")
			sche.Post(func() (ret interface{}, err error) {
				TheMgr.LoadPlayer(uid)
				return 0, nil
			})
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}
