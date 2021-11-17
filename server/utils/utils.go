package utils

import (
	"fmt"
	
	"github.com/dfklegend/cell/server/serialize"
	"github.com/dfklegend/cell/server/serialize/json"
	e "github.com/dfklegend/cell/server/errors"
)

// SerializeOrRaw serializes the interface if its not an array of bytes already
func SerializeOrRaw(serializer serialize.Serializer, v interface{}) ([]byte, error) {
	if data, ok := v.([]byte); ok {
		return data, nil
	}
	data, err := serializer.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetErrorFromPayload gets the error from payload
func GetErrorFromPayload(serializer serialize.Serializer, payload []byte) error {
	err := &e.Error{Code: e.ErrUnknownCode}
	switch serializer.(type) {
	case *json.Serializer:
		_ = serializer.Unmarshal(payload, err)
	// case *protobuf.Serializer:
	// 	pErr := &protos.Error{Code: e.ErrUnknownCode}
	// 	_ = serializer.Unmarshal(payload, pErr)
	// 	err = &e.Error{Code: pErr.Code, Message: pErr.Msg, Metadata: pErr.Metadata}
	}
	return err
}

// GetErrorPayload creates and serializes an error payload
func GetErrorPayload(serializer serialize.Serializer, err error) ([]byte, error) {
	code := e.ErrUnknownCode
	msg := err.Error()
	metadata := map[string]string{}
	if val, ok := err.(*e.Error); ok {
		code = val.Code
		metadata = val.Metadata
	}
	// errPayload := &protos.Error{
	// 	Code: code,
	// 	Msg:  msg,
	// }
	// if len(metadata) > 0 {
	// 	errPayload.Metadata = metadata
	// }
	fmt.Printf("%v%v%v", code, msg, metadata)
	errPayload := []byte("some")
	return SerializeOrRaw(serializer, errPayload)
}