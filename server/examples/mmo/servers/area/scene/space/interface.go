package space

type EntityID uint32

// 根据需求返回对象列表
type IValidator interface {
    Valid(id EntityID, dist float32) bool
    MakeResult([]EntityID) []EntityID
}

type IEntitySpace interface {
    // 添加实体
    AddEntity(entityId EntityID, x, y, z float32)
    // 移除实体
    RemoveEntity(entityId EntityID)
    // 更新实体位置
    UpdateEntityPos(entityId EntityID, x, y, z float32)

    // 在范围内搜索目标
    /**
     * 
     * @param x,y,z{float32} 中心点位置
     * @param radius{float32} 半径
     * @param validator{IValidator} 对象有效性判定器
     * @return {[]EntityID} 结果列表
    */
    SearchEntitiesInRange(x, y, z, radius float32, validator IValidator) []EntityID
}