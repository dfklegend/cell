package scene

import (        
    "log"
    "testing"     
    //"github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space/validators"      
)

// 设置两个怪物，搜索目标
func Test1_Search(t *testing.T) {
    s := NewScene()
    id1 := s.newMonster(0,0)
    id2 := s.newMonster(0,1)
    m1 := s.GetMonster(id1)
    m2 := s.GetMonster(id2)

    found := s.FindTarget(m1)
    log.Printf("%v found:%v", m1.Id, found.Id)

    found = s.FindTarget(m2)
    log.Printf("%v found:%v", m2.Id, found.Id)
}