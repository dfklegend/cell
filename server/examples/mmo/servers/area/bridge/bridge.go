package bridge

import (
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space"
)

// 解决交叉引用
type IScene interface {
    Update()
}

type IValidatorFactory interface {
    NewSimpleValidator(s IScene, id uint32) space.IValidator
}


var (
    validatorFactory IValidatorFactory
)

func SetValidatorFactory(f IValidatorFactory) {
    validatorFactory = f
}

func GetValidatorFactory() IValidatorFactory {
    return validatorFactory
}