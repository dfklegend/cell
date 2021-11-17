package test

import (
	"log"
	"sync"
	"testing"
	"time"
)

func Test_CodeTest(t *testing.T) {
	CheckUint32()
}

func SameRoutineLock() {
	mutex := sync.Mutex{}

	log.Printf("begin")
	mutex.Lock()
	log.Printf("1")
	mutex.Unlock()
	mutex.Unlock()
	mutex.Unlock()
	mutex.Lock()
	log.Printf("2")
	mutex.Lock()
	log.Printf("3")

	log.Printf("ok")
}

// Ticker.C是容量为1的channel
func CheckTicker() {
	t := time.NewTicker(1 * time.Millisecond)
	time.Sleep(10 * time.Second)
	log.Printf("%v", len(t.C))
}

func CheckUint32() {
	var a uint32
	a = uint32(0xFFFFFFFF) - 1

	for i := 0; i < 10; i++ {
		log.Printf("%v\n", a)
		a++
	}
}
