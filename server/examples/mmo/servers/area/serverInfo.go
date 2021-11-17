package area

import (    
    "sync"

    //"github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/utils/logger"
)

type AvgCalculator struct {
    total float32
    num int
}

func (self *AvgCalculator) Reset() {
    self.total = 0
    self.num = 0
}

func (self *AvgCalculator) AddValue(v float32) {
    self.total += v
    self.num ++
}

func (self *AvgCalculator) GetAvg() float32 {
    if self.num == 0 {
        return 0
    }
    return self.total/float32(self.num)
}

var serverInfo = NewServerInfo()

func GetServerInfo() *ServerInfo {
    return serverInfo
}

type ServerInfo struct {
    totalScene int
    totalMonster int
    // 平均
    monsterAttackInterval *AvgCalculator
    mutex sync.Mutex
}

func NewServerInfo() *ServerInfo {
    return &ServerInfo{
        monsterAttackInterval: &AvgCalculator{},
    }
}

func (self *ServerInfo) SetTotalScene(v int) {
    self.totalScene = v
}

func (self *ServerInfo) SetTotalMonster(v int) {
    self.totalMonster = v
}

func (self *ServerInfo) AddMonsterAttack(v float32) {
    self.mutex.Lock()
    defer self.mutex.Unlock()
    self.monsterAttackInterval.AddValue(v)
}

func (self *ServerInfo) ResetPeriodCouter() {
    self.mutex.Lock()
    defer self.mutex.Unlock()
    self.monsterAttackInterval.Reset()
}

func (self *ServerInfo) DumpInfo() {
    l := logger.Log

    l.Debugf("scenes:%v monsters:%v", self.totalScene, self.totalMonster)
    l.Debugf("monsterAttackInterval:%v", self.monsterAttackInterval.GetAvg())
}