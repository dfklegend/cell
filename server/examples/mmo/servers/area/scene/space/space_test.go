package space

import (        
    "log"
    "testing"     
    //"github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space/validators"       
)

func testPosInt(s *Space, v float32) {
    log.Printf("%v posInt:%v", v, s.getPosInt(v))
}

func testBlockId(s *Space, x, z float32) {
    log.Printf("(%v,%v) block:%v", x, z, s.calcBlockId(x, z))
}

func dumpSpace(s *Space) {
    log.Printf("width:%v", s.width)
}

func Test1_PosInt(t *testing.T) {
    s := NewSpace(-100, 100, 5)
    dumpSpace(s)

    testPosInt(s, -101)
    testPosInt(s, -100)
    testPosInt(s, -0)
    testPosInt(s, -2.0)
    testPosInt(s, 2.0)
    testPosInt(s, 100)
    testPosInt(s, 200)

    testBlockId(s, -120, -120)
    testBlockId(s, -95, -120)
    testBlockId(s, 0, 0)
    testBlockId(s, 50, 50)
}

func dumpEntity(s *Space, entityId EntityID) {
    e := s.getEntity(entityId)
    if e == nil {
        log.Printf("entity:%v is null", entityId)
        return
    }
    log.Printf("entity:%v pos(%v,%v) block:%v",
        entityId, e.PosX, e.PosZ, e.BlockId)
}

func dumpBlock(s *Space, blockId int32) {
    b := s.getBlock(blockId, false)
    if b == nil {
        log.Printf("block:%v is null", blockId)
        return
    }
    log.Printf("block:%v len:%v", blockId, len(b.Units))
    log.Printf("  %+v", b.Units)
}

func Test1_Entity(t *testing.T) {
    s := NewSpace(-100, 100, 5)

    s.AddEntity(1, 0, 0, 0)
    dumpEntity(s, 1)
    dumpBlock(s, s.calcBlockId(0, 0))

    s.UpdateEntityPos(100, 0, 0, 0)
    dumpEntity(s, 100)

    s.UpdateEntityPos(1, 100, 0, 0)
    dumpEntity(s, 1)
    dumpBlock(s, s.calcBlockId(0, 0))
    dumpBlock(s, s.calcBlockId(100, 0))


    s.AddEntity(2, 0, 0, 0)
    s.AddEntity(3, 100, 0, 0)

    // ret := s.SearchEntitiesInRange(0, 0, 0, 3, validators.NewSimpleValidator(nil, 1))
    // log.Printf("ret:%v", ret)
}