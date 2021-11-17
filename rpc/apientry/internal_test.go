package apientry

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type InMsg struct {
	Abc string `json:"abc"`
}

type OutMsg struct {
	Def string `json:"def"`
}

type Entry1 struct {
	APIEntry
}

func (self *Entry1) Join(msg *InMsg, cbFunc HandlerCBFunc) error {
	log.Printf("in join\n")

	CheckInvokeCBFunc(cbFunc, nil, &OutMsg{"dddd"})
	return nil
}

type Entry2 struct {
	APIEntry
}

func (self *Entry2) Join(d *DummySession, msg *InMsg, cbFunc HandlerCBFunc) error {
	log.Printf("in join\n")

	CheckInvokeCBFunc(cbFunc, nil, &OutMsg{"dddd"})
	return nil
}

func Test_CreateContainer(t *testing.T) {
	entry1 := NewContainer(&Entry1{}, WithName("hello"))
	entry1.ExtractHandler(defaultFormater)

	assert.Equal(t, false, entry1.HasMethod("join"), "error name")

	entry1 = NewContainer(&Entry1{}, WithNameFunc(strings.ToLower))
	entry1.ExtractHandler(defaultFormater)

	assert.Equal(t, true, entry1.HasMethod("join"), "error name")

}

func Test_CreateContainer2(t *testing.T) {
	// entry2 := NewContainer(&Entry2{}, WithNameFunc(strings.ToLower))
	// entry2.ExtractHandler(&HandlerFormater{})

	// assert.Equal(t, true, entry2.HasMethod("join"), "error name")
}

func Test_CreateCollection(t *testing.T) {
	col := NewCollection()
	col.Register(&Entry1{}, WithName("hello"))
	col.Register(&Entry1{}, WithNameFunc(strings.ToLower))
	col.Build()

	assert.Equal(t, true, col.HasMethod("hello.Join"), "error name")
	assert.Equal(t, true, col.HasMethod("entry1.join"), "error name")
}

func Test_CallRPC(t *testing.T) {
	col := NewCollection()
	col.Register(&Entry1{}, WithName("hello"))
	col.Register(&Entry1{}, WithNameFunc(strings.ToLower))
	col.Build()

	col.Call("hello.Join", []byte("ddd"), func(e error, result interface{}) {
		log.Printf("got result:%v\n", string(result.([]byte)))
	}, nil)
	col.Call("hello.Join", []byte("ddd"), nil, nil)
}
