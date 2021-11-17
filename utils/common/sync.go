package common


type IMutex interface {
    Lock()
    Unlock()
    RLock()
    RUnlock()
}

// 提供空的锁，便于测试同步(替代掉sync.Mutex,sync.RWMutex)
type FakeMutex struct {    
}

func (self *FakeMutex) Lock() {}
func (self *FakeMutex) Unlock() {}
func (self *FakeMutex) RLock() {}
func (self *FakeMutex) RUnlock() {}