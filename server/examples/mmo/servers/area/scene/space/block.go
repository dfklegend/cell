package space

import (
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/utils"
)

type Block struct {
    //
    Id      int32
    // 对象列表
    Units   map[EntityID]*Entity
}

func NewBlock(id int32) *Block {
    return &Block{
        Id: id,
        Units: make(map[EntityID]*Entity),
    }
}

func (self *Block) Add(e *Entity) {
    self.Units[e.Id] = e
}

func (self *Block) Del(id EntityID) {
    delete(self.Units, id)
}

func (self *Block) SearchEntitiesInRange(x, y, z, radius float32, validator IValidator) []EntityID {
    if len(self.Units) == 0 {
        return nil
    }    

    ret := make([]EntityID, 0, 16)
    for k, e := range(self.Units) {
        dist := utils.CalcDistXZ(e.PosX, e.PosZ, x, z)
        if dist > radius {
            continue
        }
        if !validator.Valid(e.Id, dist) {
            continue
        }
        ret = append(ret, k)
    }
    return ret
}

