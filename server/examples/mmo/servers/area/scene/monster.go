package scene

import (    
    "math"
    "math/rand"

    "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/server/examples/mmo/servers/area"
)

// 怪物
type Monster struct {
    Id uint32
    scene *Scene
    HP int
    // 位置
    PosX, PosY float32
    nextAttack int64
    lastAttack int64
    nextMove int64
    corpseDisappearTime int64
}

func NewMonster(id uint32) *Monster {    
    return &Monster{
        Id: id,
        HP: 100,
        nextAttack: common.NowMs() + rand.Int63n(1000),
    }
}

func (self* Monster) IsDead() bool {
    return self.HP == 0
}

func (self* Monster) IsOver() bool {
    return self.IsDead() && common.NowMs() >= self.corpseDisappearTime
}

func (self* Monster) SetScene(scene *Scene) {
    self.scene = scene
}


// 每隔xs，随机移动一格
// 每秒，周边找个目标K一下
func (self *Monster) Update(s *Scene) {
    self.updateAttack(s)
    self.updateMove()
}

func (self *Monster) updateMove() {
    now := common.NowMs()
    if now < self.nextMove {
        return
    } 
    self.nextMove = now

    self.PosX += common.RandFloat32(-1.0, 1.0)    
    self.PosX = self.clamp(self.PosX, MinX, MaxX)
    self.PosY += common.RandFloat32(-1.0, 1.0)
    self.PosY = self.clamp(self.PosY, MinY, MaxY)
    self.scene.UpdateMonsterPos(self)
}

func (self *Monster) clamp(pos, min, max float32) float32 {
    if pos > max {
        return max
    }
    if pos < min {
        return min
    }
    return pos
}

func (self *Monster) updateAttack(s *Scene) {
    now := common.NowMs()
    if now < self.nextAttack {
        return
    } 
    if self.lastAttack > 0 {
        off := now - self.lastAttack
        area.GetServerInfo().AddMonsterAttack(float32(off))
    }
    self.lastAttack = now

    self.nextAttack = now + 1000   

    for i := 0; i < 10; i ++ {
        s.FindTarget(self)
    } 
    tar := s.FindTarget(self)
    if tar != nil {
        self.Attack(tar)        
    }    
}

func (self *Monster) Attack(tar *Monster) {    
    if tar == nil || tar.IsDead() {
        return
    }
    tar.HP --
    if tar.HP < 0 {
        tar.HP = 0
        tar.OnDead()
    }    
}

func (self *Monster) DistTo(tar *Monster) float32 {
    off1 := self.PosX - tar.PosX
    off2 := self.PosY - tar.PosY
    return float32(math.Sqrt(float64(off1*off1 + off2*off2)) )
}

func (self *Monster) OnDead() {
    self.corpseDisappearTime = common.NowMs() + 10*1000
}


