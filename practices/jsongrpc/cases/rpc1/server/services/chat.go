
package services

import (
	"fmt"

	"dfk.com/practices/jsongrpc/libs/rpc/service"	
	"dfk.com/practices/jsongrpc/cases/rpc1/protos"
)

type Chat struct {
	service.Base
}

func (rs *Chat) Join(dummy int, msg *protos.InMsg) error {
	fmt.Printf("Chat.Join:%v\n", msg)
	return nil
}