package services

import (	
	"dfk.com/practices/jsongrpc/libs/rpc/service"
	"dfk.com/practices/jsongrpc/cases/rpc1/protos"
)

type Some struct {
	service.Base
}

func (rs *Some) JoinRoom(dummy int, msg *protos.InMsg) error {
	return nil
}