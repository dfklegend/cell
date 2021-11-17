package validators

import (
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/bridge"
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space"
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene"
)


type SimpleValidator struct {
    curScene    bridge.IScene
    src         space.EntityID
    tarId       space.EntityID
    curDist     float32
}

func NewSimpleValidator(s bridge.IScene, srcId uint32) *SimpleValidator {
    v := &SimpleValidator {
        curScene: s,
        src: space.EntityID(srcId),
        curDist: 0.0,
    }
    return v
}

func (self *SimpleValidator) Valid(id space.EntityID, dist float32) bool {
    if self.src == id {
        return false
    }

    // 检查是否死亡
    if self.curScene != nil {
        s := self.curScene.(*scene.Scene)
        m := s.GetMonster(uint32(id))
        if m != nil && m.IsDead() {
            return false
        }    
    }    

    if self.tarId == 0 || dist < self.curDist {
        self.tarId = id
        self.curDist = dist
        return true   
    }    
    return false
}

func (self *SimpleValidator) MakeResult([]space.EntityID) []space.EntityID {
    ret := make([]space.EntityID, 0, 1)
    ret = append(ret, self.tarId)
    return ret
}