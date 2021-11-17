package channel

import (
	"log"	
	"testing"	
	"fmt"
	"time"
	"math/rand"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell/server/common"
)


func Test_Group(t *testing.T) {
	g := NewFrontGroup()

	for i := 0; i < 256; i ++ {
		g.NetIds = append(g.NetIds, common.NetIdType(i) )
	}

	log.Printf("%+v\n %v %v", g.NetIds, len(g.NetIds), cap(g.NetIds))
}


func Test_Remove(t *testing.T) {
	g := NewFrontGroup()

	for i := 0; i < 100; i ++ {
		g.NetIds = append(g.NetIds, common.NetIdType(i) )
	}

	g.Remove(0)
	log.Printf("%+v\n%v %v", g.NetIds, len(g.NetIds), cap(g.NetIds))
	// 期望元素里没有0,且长度为99
	
	assert.Equal(t, true, g.FindIndex(0) == -1 && len(g.NetIds) == 99,
	 	"mismatch")

	g.Remove(99)
	log.Printf("%+v\n%v %v", g.NetIds, len(g.NetIds), cap(g.NetIds))
	assert.Equal(t, true, g.FindIndex(99) == -1 && len(g.NetIds) == 98,
	 	"mismatch")
	// 期望元素里没99
	g.Remove(common.NetIdType(50))
	log.Printf("%+v\n%v %v", g.NetIds, len(g.NetIds), cap(g.NetIds))
	assert.Equal(t, true, g.FindIndex(50) == -1 && len(g.NetIds) == 97,
	 	"mismatch")
}


// 并发测试
func Test_Concurrence(t *testing.T) {
	log.Println("---- Test_Concurrence ----")
	run := true

	s := channelService
	// create service
	for i := 0; i < 10; i ++ {
		go func() {
			for run {
				c := s.AddChannel(fmt.Sprintf("c%v", rand.Intn(100)))
				c.Add(fmt.Sprintf("f%v", rand.Intn(100)), common.NetIdType(rand.Intn(100)) )
			}
		}()	
	}	

	for i := 0; i < 10; i ++ {
		go func() {
			for run {
				c := s.GetChannel(fmt.Sprintf("c%v", rand.Intn(100)))
				if c != nil {
					c.Add(fmt.Sprintf("f%v", rand.Intn(100)), common.NetIdType(rand.Intn(100)) )	

					// c.Range(func( k, v interface{}) bool {
					// 	n := k.(string)
				 //        g := v.(*FrontGroup)

				 //        g.Lock()
				 //        log.Printf("c:%v:%v\n", n, g.NetIds)
				 //        g.Unlock()
				 //        return true
					// })
				}			
			}
		}()
	}	

	time.Sleep(3 * time.Second)
	run = false

	// dump result
	for i := 0; i < 100; i ++ {
		channelName := fmt.Sprintf("c%v", i)
		c := s.GetChannel(channelName)
		if c != nil {
			log.Printf("channel:%v\n", channelName)
			c.Range(func( k, v interface{}) bool {
				n := k.(string)
		        g := v.(*FrontGroup)

		        g.Lock()
		        log.Printf("%v:%v\n", n, g.NetIds)
		        g.Unlock()
		        return true
			})
		}	
	}

	time.Sleep(1 * time.Second)
}
