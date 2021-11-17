package sche

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/dfklegend/cell/utils/common"
)

var (
	TheScheMgr     *ScheMgr = NewScheMgr()
	DefaultScheMgr *ScheMgr = TheScheMgr
)

func init() {
	//fmt.Println("sche.Init")
}

// 最简单的匿名函数
type CBFunc func() (interface{}, error)

// --------------------------

type RunTask struct {
	id uint32
	cb CBFunc
	// 返回结果
	chanRet chan interface{}
	err     error
}

func (s *RunTask) Done(r interface{}, e error) {
	s.err = e
	s.chanRet <- r
}

// 防止死锁
func (s *RunTask) WaitDone() (interface{}, error) {
	const timeout = 99
	select {
	case r := <-s.chanRet:
		return r, s.err
	case <-time.After(timeout * time.Second):
		//log.Panicf("致命错误: 是否向本sche投递并waitDone:%v", s)
		panic(fmt.Sprintf("致命错误: 是否向本sche投递并waitDone:%v", s))
		//return nil, nil
	}
}

// --------------------------
type RunTaskIdService struct {
	nextId uint32
}

func (s *RunTaskIdService) AllocId() uint32 {
	return atomic.AddUint32(&s.nextId, 1)
}

var (
	s_runTaskIdService RunTaskIdService = RunTaskIdService{
		nextId: 1,
	}
)

// --------------------------

type ISche interface {
	AddTask(cb CBFunc) *RunTask
	Wait()
}

type Sche struct {
	chanTask  chan *RunTask
	chanClose chan int
}

const (	
	ScheQueueSize = 9999
)

func NewSche() *Sche {
	return &Sche{
		chanTask:  make(chan *RunTask, ScheQueueSize),
		chanClose: make(chan int, 1),
	}
}

func (s *Sche) Stop() {
	close(s.chanTask)
	close(s.chanClose)
}

// TODO: 如果投递到本routine然后又WaitDone，会死锁
// 如何解决?或者警告
// 安全考虑，如果某个
func (s *Sche) Post(cb CBFunc) *RunTask {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic in sche.Post:%v", err)	
			log.Printf(common.GetStackStr())
		}
	}()

	t := &RunTask{
		id:      s_runTaskIdService.AllocId(),
		cb:      cb,
		chanRet: make(chan interface{}, 1),
	}

	if len(s.chanTask) == ScheQueueSize {
		log.Printf("sche channel is full:%v %v", s, len(s.chanTask))
		// 阻塞
	}
	s.chanTask <- t
	return t
}

// 集成在其他流程中
/*
	select {
	case t := <- sche.GetChanTask():
		sche.DoTask(t)
	}
*/
func (s *Sche) GetChanTask() <-chan *RunTask {
	return s.chanTask
}

func (s *Sche) doTask(t *RunTask) (r interface{}, e error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("panic in task do:%v", err)
			e = errors.New(fmt.Sprintf("panic in doTask:%v", err))
		}
	}()
	
	return t.cb()
}

func (s *Sche) DoTask(t *RunTask) {
	t.Done(s.doTask(t))
}

// 便于测试，自己go sche.Handler
func (s *Sche) Handler() {
	for true {
		select {
		case t := <-s.chanTask:
			if t != nil {
				s.DoTask(t)
			}			
		case <-s.chanClose:
			return
		}
	}
}

// --------------------------
// 管理器,通过名字来创建和访问sche
type ScheMgr struct {
	sches map[string]*Sche
	mutex sync.Mutex
}

func NewScheMgr() *ScheMgr {
	return &ScheMgr{
		sches: make(map[string]*Sche),
		//mutex:
	}
}

// createIfMiss
func (s *ScheMgr) GetSche(name string) *Sche {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	one := s.sches[name]
	if one == nil {
		s.sches[name] = NewSche()
		return s.sches[name]
	}
	return one
}

func (s *ScheMgr) hasSche(name string) bool {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	one := s.sches[name]
	return one != nil
}

func (s *ScheMgr) DelSche(name string) {
	defer s.mutex.Unlock()

	s.mutex.Lock()
	delete(s.sches, name)
}
