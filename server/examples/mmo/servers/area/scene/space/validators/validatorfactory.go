package validators

import (
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/bridge"
    "github.com/dfklegend/cell/server/examples/mmo/servers/area/scene/space"
)

func init() {
    v := &ValidatorFactory{}
    bridge.SetValidatorFactory(v)
}

func Visit() {    
}

type ValidatorFactory struct {

}

func (self *ValidatorFactory) NewSimpleValidator(s bridge.IScene, id uint32) space.IValidator {
    return NewSimpleValidator(s, id)
}

