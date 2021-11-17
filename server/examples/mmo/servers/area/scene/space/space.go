package space

// ------------

type Space struct {
    min, max        float32
    width           int32
    gridSize        float32
    blocks          map[int32]*Block
    entities        map[EntityID]*Entity
}

// 一个正方形的块
// 2维，Y轴无视
func NewSpace(min, max, gridSize float32) *Space {
    v := &Space {
        min: min,
        max: max,
        gridSize: gridSize,
        blocks: make(map[int32]*Block),
        entities: make(map[EntityID]*Entity),
    }
    v.Init()
    return v
}

func (self* Space) Init() {
    self.width = int32((self.max - self.min)/self.gridSize) + 1
}

func (self *Space) clamp(v float32) float32 {
    if v < self.min {
        return self.min
    }
    if v > self.max {
        return self.max
    }
    return v
}

func (self *Space) getPosInt(v float32) int32 {
    v = self.clamp(v)
    return int32((v-self.min)/self.gridSize)
}

func (self *Space) calcBlockId(x, z float32) int32 {
    col := self.getPosInt(x)
    row := self.getPosInt(z)
    return self.getBlockIdFromRowCol(row, col)
}

func (self *Space) getBlockIdFromRowCol(row, col int32) int32 {
    return row * self.width + col
}

func (self *Space) AddEntity(entityId EntityID, x, y, z float32) {
    exsit := self.getEntity(entityId)
    if exsit != nil {
        return
    }
    e := NewEntity(entityId)
    self.entities[entityId] = e

    self.updateEntityPos(e, x, y, z)
}

func (self *Space) RemoveEntity(entityId EntityID) {
    e := self.getEntity(entityId)
    if e == nil {
        return
    }
    oldBlock := self.getBlock(e.BlockId, false)
    if oldBlock == nil {
        return
    }

    oldBlock.Del(entityId)
    delete(self.entities, entityId)
}

func (self *Space) getEntity(entityId EntityID) *Entity {
    e, _ := self.entities[entityId]
    return e
}

func (self *Space) getBlock(blockId int32, createIfMiss bool) *Block {
    b, _ := self.blocks[blockId]
    if b == nil {
        if !createIfMiss {
            return nil
        }

        // 添加一个
        nb := NewBlock(blockId)
        self.blocks[blockId] = nb
        return nb
    }
    return b
}

// 根据位置刷新block
func (self *Space) updateEntityPos(e *Entity, x, y, z float32) {
    if e == nil {
        return
    }
    newBlockId := self.calcBlockId(x, z)
    if newBlockId == e.BlockId {
        e.SetPos(x, y, z)
        return
    }
    oldBlock := self.getBlock(e.BlockId, false)
    newBlock := self.getBlock(newBlockId, true)

    if oldBlock != nil {
        oldBlock.Del(e.Id)
    }
    e.SetPos(x, y, z)
    newBlock.Add(e)
    e.BlockId = newBlockId
}

func (self *Space) UpdateEntityPos(entityId EntityID, x, y, z float32) {    
    self.updateEntityPos(self.getEntity(entityId), x, y, z)
}

func (self *Space) SearchEntitiesInRange(x, y, z, radius float32, validator IValidator) []EntityID {
    ret := make([]EntityID, 0, 16)

    xStart := self.getPosInt(x - radius)
    xEnd := self.getPosInt(x + radius)
    zStart := self.getPosInt(z - radius)
    zEnd := self.getPosInt(z + radius)
    for row := zStart; row <= zEnd; row ++ {
        for col := xStart; col <= xEnd; col ++ {
            blockId := self.getBlockIdFromRowCol(row, col)
            block := self.getBlock(blockId, false)
            if block == nil {
                continue
            }
            blockRet := block.SearchEntitiesInRange(x, y, z, radius, validator)
            if blockRet != nil && len(blockRet) > 0 {
                ret = append(ret, blockRet...)
            }
        }
    }
    if len(ret) > 0 {
        return validator.MakeResult(ret)
    }
    return ret
}