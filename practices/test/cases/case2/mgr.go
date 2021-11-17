package main

import (
	"fmt"
	"time"

	scheP "github.com/dfklegend/cell/practices/test/libs/sche"
	"github.com/dfklegend/cell/practices/test/libs/util"
	waterfall "github.com/dfklegend/cell/practices/test/libs/waterfall2"
)

var (
	TheMgr *Mgr
	TheDB  *DB
)

func init() {
	TheMgr = &Mgr{}
	TheDB = &DB{}
}

type Player struct {
	info string
}

type Mgr struct {
}

func (s *Mgr) Handle() {
	var sche = scheP.TheScheMgr.GetSche("mgr")

	for true {
		select {
		case t := <-sche.GetChanTask():
			sche.DoTask(t)
		}
	}
}

func (s *Mgr) CreatePlayer() *Player {
	return &Player{}
}

// 虚拟的载入玩家
func (s *Mgr) LoadPlayer(uid int) {
	var sche = scheP.TheScheMgr.GetSche("mgr")

	// 创建一个玩家对象
	// 请求 db 读取玩家数据
	// player对象读取玩家数据,构建完毕
	var p *Player
	fmt.Printf("LoadPlayer %d routine:%d\n", uid, util.GetRoutineID())
	waterfall.Waterfall(sche,
		[]waterfall.Task{func(args ...interface{}) {
			callback, _ := args[0].(waterfall.Callback)
			fmt.Printf("CreatePlayer %d routine:%d\n", uid, util.GetRoutineID())
			p = s.CreatePlayer()
			callback(false)
		}, func(args ...interface{}) {
			callback, _ := args[0].(waterfall.Callback)
			// 请求 db 读取玩家数据
			s.goDBLoadPlayer(func(args1 ...interface{}) {
				fmt.Printf("goDBLoadPlayer cb %d routine:%d\n", uid, util.GetRoutineID())
				playerInfo, _ := args1[0].(string)
				callback(false, playerInfo)
			})
		}, func(args ...interface{}) {
			callback, _ := args[0].(waterfall.Callback)
			playerInfo, _ := args[1].(string)
			fmt.Printf("set p.info %d routine:%d\n", uid, util.GetRoutineID())
			p.info = fmt.Sprintf("%s%d", playerInfo, uid)
			callback(false)
		}}, func(args ...interface{}) {
			fmt.Printf("final %d routine:%d\n", uid, util.GetRoutineID())
			fmt.Println(args...)
			fmt.Println("p:", p)
		})
}

// 避免阻塞Mgr routine
func (s *Mgr) goDBLoadPlayer(callback waterfall.Task) {
	var dbSche = scheP.TheScheMgr.GetSche("db")
	go func() {
		fmt.Println("goDBLoadPlayer go begin:", util.GetRoutineID())
		t := dbSche.Post(func() (ret interface{}, err error) {
			ret = TheDB.LoadPlayer()
			callback(ret)
			return ret, nil
		})
		t.WaitDone()
		fmt.Println("goDBLoadPlayer waitDone over:", util.GetRoutineID())
	}()
}

// -----------------------
type DB struct{}

func (s *DB) Handle() {
	var sche = scheP.TheScheMgr.GetSche("db")

	for true {
		select {
		case t := <-sche.GetChanTask():
			sche.DoTask(t)
		}
	}
}

// db部分 如何并发?
// 一般数据库读取也是 "同步"接口，如何处理
// 建立一个连接池，路由保证同一个id指向一个连接
func (s *DB) LoadPlayer() string {
	fmt.Println("db.LoadPlayer enter:", util.GetRoutineID())
	time.Sleep(3 * time.Second)
	fmt.Println("db.LoadPlayer leave:", util.GetRoutineID())
	return "abc"
}
