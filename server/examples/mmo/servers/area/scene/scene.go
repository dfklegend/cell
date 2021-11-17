package scene

import (    
    //"log"    
    "math/rand"
    "fmt"
    "time"      

    "github.com/dfklegend/cell/utils/common"    
    "github.com/dfklegend/cell/utils/runservice"
    "github.com/dfklegend/cell/utils/timer"

    spacep "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space"    
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/bridge"    
    "github.com/dfklegend/cell/utils/logger"    
)

type Scene struct  {
    sceneId uint32
    monsters map[uint32]*Monster
    over bool
    maxMonsters int
    nextCanSpawn int64

    idService *common.SerialIdService
    // 是否有独立的runService
    runService *runservice.StandardRunService
    externRunService bool
    running bool

    space spacep.IEntitySpace

    timerCheckDead timer.TimerIdType
    eventTestId uint64
}

func NewScene(id uint32) *Scene {
    
    space := spacep.NewSpace(MinX, MaxX, 5.0)

    return &Scene {
        sceneId: id,
        monsters: make(map[uint32]*Monster),
        over: false,
        maxMonsters: 100,
        idService: common.NewSerialIdService(),
        running: true,
        space: space,
        externRunService: false,
    }
}

func (self *Scene) SetRunService(v *runservice.StandardRunService) {
    self.runService = v
    self.externRunService = true
}

func (self *Scene) Update() {
    for k, v := range(self.monsters) {
        v.Update(self)
        if v.IsOver() {
            delete(self.monsters, k)
            self.space.RemoveEntity(spacep.EntityID(k))
            //logger.Log.Debugf("delete monster")            
        }
    }
    self.updateSpawn()    
}

func (self *Scene) newMonster(x, y float32) uint32 {
    id := self.idService.AllocId()

    monster := NewMonster(id)
    monster.PosX = x
    monster.PosY = y
    monster.SetScene(self)
    
    self.monsters[id] = monster  
    self.space.AddEntity(spacep.EntityID(id), monster.PosX, 0, monster.PosY)  
    return id
}

func (self *Scene) updateSpawn() {
    // 尝试刷怪
    if len(self.monsters) >= self.maxMonsters {
        return
    }

    now := common.NowMs()
    if now < self.nextCanSpawn {
        return
    }
    self.nextCanSpawn = now + 100

    // new monster
    for i := 0; i < 100; i ++ {
        self.newMonster(common.RandFloat32(MinX, MaxX),
            common.RandFloat32(MinY, MaxY))        
    }    

    //logger.Log.Debugf("new monster")    
}

func (self *Scene) SetOver() {
    self.over = true
}

func (self *Scene) IsOver() bool {
    return self.over
}

// 索敌
// 根据位置找个最近的敌人
func (self *Scene) FindTarget(src *Monster) *Monster {
    //return src
    //return self.FindTargetSimple(src)
    return self.FindTargetSpace(src)
}

func (self *Scene) FindTargetSimple(src *Monster) *Monster {

    var found *Monster
    var dist float32
    for _, m := range(self.monsters) {
        if m == src {
            continue
        }
        if found == nil {
            found = m
            dist = src.DistTo(m)
            continue
        }

        newDist := src.DistTo(m)
        if newDist < dist {
            found = m
            dist = newDist
        }
    }
    return found
}

func (self *Scene) FindTargetSpace(src *Monster) *Monster {    
    ids := self.space.SearchEntitiesInRange(src.PosX, 0, src.PosY, 6.0,
        bridge.GetValidatorFactory().NewSimpleValidator(self, src.Id))
    if ids != nil && len(ids) > 0 {
        return self.GetMonster(uint32(ids[0]))
    }
    return nil
}

func (self *Scene) UpdateMonsterPos(m *Monster) {
    self.space.UpdateEntityPos(spacep.EntityID(m.Id), 
        m.PosX, 0, m.PosY)
}

func (self *Scene) GetMonsterNum() int {
    return len(self.monsters)
}

func (self *Scene) GetMonster(id uint32) *Monster {
    m, _ := self.monsters[id]
    return m
}

func (self *Scene) Start() {
    self.startLogic()
}

// routine里跑
func (self *Scene) _StartGo() {
    //logger.Log.Debugf("---- StartGo ----")        

    // start runService
    name := fmt.Sprintf("scene:%v", self.sceneId)
    self.runService = runservice.NewStandardRunService(name)

    self.runService.GetTimerMgr().AddTimer(30*time.Millisecond, 
        func(args ...interface{}) {            
        self.Update()            
    })     
    self.runService.Start()
    self.startLogic()
}

func (self *Scene) startLogic() {
    // 注册update
    self.runService.GetTimerMgr().AddTimer(30*time.Millisecond, 
        func(args ...interface{}) {            
        self.Update()            
    })   

    // 测试事件
    self.eventTestId = self.runService.GetEventCenter().GSubscribe("30s", func(args ...interface{}){
        logger.Log.Debugf("30s tirggle %v in:%v args:%v",
            self.sceneId, common.GetRoutineID(), args)
    })   
    
    self.timerCheckDead = self.runService.GetTimerMgr().AddTimer(1*time.Second, 
        func(args ...interface{}) {            
        // 概率死亡
        if rand.Float32() <= 0.01 {
            logger.Log.Debugf("scene:%v SetOver", self.sceneId)
            self.SetOver()
        }
    })   
}

func (self *Scene) Stop() {
    self.running = false    

    if !self.externRunService {
        self.runService.Stop()
        self.timerCheckDead = 0
        self.eventTestId = 0
    } else {
        self.runService.GetTimerMgr().Cancel(self.timerCheckDead)
        self.timerCheckDead = 0

        self.runService.GetEventCenter().GUnsubscribe("30s", self.eventTestId)
        self.eventTestId = 0    
    }
}

