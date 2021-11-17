package event

import (
	"log"
	"testing"
	"time"
	"sync"
)

// 测试timer基本流程
// 创建,cancel
func Test_Normal(t *testing.T) {
	center := NewLocalEventCenter(false)

	id1 := center.Subscribe("event1", func(args ...interface{}) {
		log.Println(args)
	}, 1, "haha")
	id1_2 := center.Subscribe("event1", func(args ...interface{}) {
		log.Println(args)
	})

	center.Publish("event1", 3, 4)
	center.Publish("event1", 5, 6)
	center.Unsubscribe("event1", id1)
	center.Publish("event1")

	log.Printf("ids:%v %v", id1, id1_2)
}

func localRun(index int) {
	center := NewLocalEventCenter(false)
	center.SetLocalUseChan(true)

	id1 := center.Subscribe("eventLocal", func(args ...interface{}) {
		log.Printf("%v eventLocal %v", index, args)
	}, "selfArg1")

	id2 := center.GSubscribe("event2", func(args ...interface{}) {
		log.Printf("%v event2 %v", index, args)
	}, "selfArg2")

	running := true
	go func() {
		time.Sleep(5 * time.Second)
		running = false
	}()

	go func() {
		time.Sleep(2 * time.Second)
		center.Publish("eventLocal", "from local")		
		center.Publish("event2", "from local")		
	}()

	for running {
		select {
		case e := <-center.GetChanEvent():
			center.DoEvent(e)
		default:
		}
	}	

	log.Printf("running:%v", running)
	log.Printf("ids:%v %v", id1, id2)
	center.Clear()
}

func Test_Global(t *testing.T) {

	go func() {
		time.Sleep(1 * time.Second)
		GetGlobalEC().Publish("eventLocal", "fromGlobal")
		GetGlobalEC().Publish("event2", "fromGlobal")
		time.Sleep(100 * time.Millisecond)		
	}()
	
	go localRun(1)
	go localRun(2)

	time.Sleep(3 * time.Second)
	GetGlobalEC().DumpInfo()
	time.Sleep(4 * time.Second)
	GetGlobalEC().DumpInfo()
}

// 测试并发安全性
// 并发订阅/取消订阅，publish
func doTest_Concurrence(t *testing.T, fakeMutex bool) {	
	useFakeMutex = fakeMutex
	center := NewLocalEventCenter(false)

	running := true
	ids := make([]uint64, 9999)
	var mutex sync.Mutex
	for i := 0; i < 3; i ++ {
		go func() {			
			for running {
				id := center.Subscribe("event1", func(args ...interface{}) {
					//log.Printf("1")
				})
				mutex.Lock()
				ids = append(ids, id)
				mutex.Unlock()
			}
		}()	
	}

	for i := 0; i < 1; i ++ {
		go func() {			
			for running {
				mutex.Lock()
				if len(ids) > 0 {
					id := ids[0]
					ids = ids[1:]

					center.Unsubscribe("event1", id)
				}
				mutex.Unlock()	
			}
		}()	
	}


	for i := 0; i < 3; i ++ {
		go func() {			
			for running {
				center.Publish("event1",1)
			}
		}()	
	}

	go func() {
		t := time.NewTicker(1*time.Second)

		for running {
			select {
			case <-t.C:
				center.DumpInfo("event1")
			}
		}
	}()

	time.Sleep(5*time.Second)	
	running = false
}

func Test_ConcurrenceNormal(t *testing.T) {
	log.Printf("---- Test_ConcurrenceNormal ----")
	doTest_Concurrence(t, false)
}

func Test_ConcurrenceFakeMutex(t *testing.T) {
	log.Printf("---- Test_ConcurrenceFakeMutex ----")
	log.Printf("---- 应该有map竞争冲突 ----")
	doTest_Concurrence(t, true)
}