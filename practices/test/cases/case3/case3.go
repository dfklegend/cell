package main
 
import (
	"sync"
	"time"
	"fmt"
)
 
var lock sync.RWMutex
 
func main() {
	go func() {
		outRLock()
	}()
	time.Sleep(time.Second)
	wLock()


	time.Sleep(time.Second*10)
}
 
func outRLock() {
	fmt.Println("outRLock in")
	lock.RLock()
	defer lock.RUnlock()
	time.Sleep(time.Second * 2)
	innerRLock()
	fmt.Println("outRLock out")
}
 
func innerRLock() {
	fmt.Println("innerRLock in")
	lock.RLock()
	fmt.Println("innerRLock out")
	defer lock.RUnlock()
}
 
func wLock() {
	fmt.Println("wlock in")
	lock.Lock()	
	fmt.Println("wlock out")
	defer lock.Unlock()
}