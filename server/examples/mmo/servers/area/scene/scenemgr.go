package scene

import (        
    "sync"
    "time"    
    "reflect"
    "math/rand"

    "github.com/dfklegend/cell/utils/common"
    "github.com/dfklegend/cell/utils/runservice"
    "github.com/dfklegend/cell/utils/timer"
    "github.com/dfklegend/cell/utils/event"
    "github.com/dfklegend/cell/utils/sche"    
    "github.com/dfklegend/cell/utils/logger"    

    "github.com/dfklegend/cell/server/examples/mmo/servers/area"
)

var sceneRoutineNum = 5

var (
    sceneIdService = common.NewSerialIdService()
    sceneMgr = newSceneMgr()
)

func SetAllSceneOnRoutine(v bool) {
    if(v) {
        sceneRoutineNum = 1    
    }    
}

// 管理
type SceneMgr struct {
    // id:*Scene
    scenes sync.Map
    lastSize int
    runService *runservice.StandardRunService
    runPool *RunPool

    timerMgr *timer.TimerMgr

    tickerUpdate *time.Ticker
    maxScene int

    nextSetInfo int64

    timerDumpInfo timer.TimerIdType
    eventTestIndex int
}

func newSceneMgr() *SceneMgr {
    return &SceneMgr{
        runPool: NewRunPool("scenepool"), 
        runService: runservice.NewStandardRunService("sceneMgr"),        
        maxScene: 100,        
    }
}

func GetSceneMgr() *SceneMgr {
    return sceneMgr
}

func (self *SceneMgr) allocSceneId() uint32 {
    return sceneIdService.AllocId()
}

func (self *SceneMgr) selectRunService() *runservice.StandardRunService {
    return self.runPool.GetService(rand.Intn(5))
    //return self.runService
}

func (self *SceneMgr) CreateScene() uint32 {
    
    id := self.allocSceneId()
    s := NewScene(id)    
    self.scenes.Store(id, s)

    logger.Log.Debugf("CreateScene:%v", id)
    // if allSceneOneRoutine {
    //     s.SetExternRunService(self.selectRunService()) 
    //     s.Start()
    // } else {
    //     s.StartGo()    
    // }

    s.SetRunService(self.selectRunService()) 
    s.Start()
    return id
}

func (self *SceneMgr) Start() {
    self.timerMgr = self.runService.GetTimerMgr()

    self.runService.Start()
    self.startUpdate()    

    self.timerDumpInfo = self.timerMgr.AddTimer(5*time.Second, self.dumpInfo)
    logger.Log.Debugf("SceneMgr.Start")

    self.timerMgr.AddTimer(6*time.Second, self.eventTest)
    self.timerMgr.AddTimer(1*time.Second, self.updateSceneDead)

    self.runPool.Start(sceneRoutineNum)
}

func (self *SceneMgr) dumpInfo(args ...interface{}) {
    info := area.GetServerInfo()
    info.DumpInfo()
    info.ResetPeriodCouter()
}

func (self *SceneMgr) eventTest(args ...interface{}) {
    event.GetGlobalEC().Publish("30s", self.eventTestIndex)
    self.eventTestIndex ++
}

func (self *SceneMgr) Stop() {
    self.timerMgr.Cancel(self.timerDumpInfo)
    self.timerDumpInfo = 0

    self.runService.Stop()
    self.stopUpdate()
}

func (self *SceneMgr) startUpdate() {
    t := time.NewTicker(30*time.Millisecond)
    self.tickerUpdate = t
    self.runService.GetSelector().AddSelector(
        sche.NewFuncSelector(reflect.ValueOf(t.C),
            func(v reflect.Value, recvOk bool) {                
                self.update()            
    }))
}

func (self *SceneMgr) stopUpdate() {
    self.tickerUpdate.Stop()
}

func (self* SceneMgr) update() {
    // if(allSceneOneRoutine) {
    //     self.updateScenes()    
    // } else {
    //     self.calcSceneLen()
    // }
    self.calcSceneLen()
    self.updateSpawnScene()
    self.updateInfo()
}

func (self* SceneMgr) updateScenes() {
    size := 0
    self.scenes.Range(func( k, v interface{}) bool {
        s := v.(*Scene)
        s.Update()
        size ++
        return true
    })
    self.lastSize = size    
}

func (self* SceneMgr) calcSceneLen() {
    size := 0
    self.scenes.Range(func( k, v interface{}) bool {       
        size ++
        return true
    })
    self.lastSize = size    
}

func (self* SceneMgr) updateSpawnScene() {
    if self.lastSize >= self.maxScene {
        return
    }
    self.CreateScene()
}

func (self* SceneMgr) updateSceneDead(args ...interface{}) {
    // 判断scene死亡
    self.scenes.Range(func( k, v interface{}) bool {       
        theScene := v.(*Scene)
        if theScene.IsOver() {            
            logger.Log.Debugf("scene:%v Stop", k)
            theScene.Stop()
            self.scenes.Delete(k)            
        }
        return true
    })
}

// 更新数据
func (self* SceneMgr) updateInfo() { 
    // 每秒更新下  
    now := common.NowMs()
    if now < self.nextSetInfo {
        return
    }
    self.nextSetInfo = now + 1000
    info := area.GetServerInfo()

    info.SetTotalScene(self.lastSize)
    // 统计每个场景的怪
    num := 0
    self.scenes.Range(func( k, v interface{}) bool {       
        s := v.(*Scene)
        num += s.GetMonsterNum()
        return true
    })
    info.SetTotalMonster(num)
}