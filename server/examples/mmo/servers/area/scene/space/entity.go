package space

// 结构
type Entity struct {
    // 位置
    PosX, PosY, PosZ float32
    // 对象id
    Id EntityID
    // 当前blockId
    BlockId int32
}

func NewEntity(id EntityID) *Entity {
    return &Entity{
        Id: id,
        BlockId: -1,
    }
}

func (self *Entity) SetPos(x, y, z float32) {
    self.PosX = x
    self.PosY = y
    self.PosZ = z
}