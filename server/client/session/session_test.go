package session

import (
	"log"	
	"testing"	
)


func Test_BackSession(t *testing.T) {
	//var is IServerSession

	f := NewFrontSession()
	f.Set("k1", 1)
	f.Set("k2", "ffff")
	
	log.Printf("%v %v\n", f.Get("k1", 0), f.Get("k2", ""))

	str := f.ToJson()
	b := NewBackSession(str)	
	//is = b
	log.Printf("%+v %v\n", b, b.ToJson())
	b.Set("k3", 2)
	b.Set("k4", 2)
	log.Printf("%+v %v\n", b, b.ToJson())
}

func TestCB(t *testing.T) {
	mgr := GetFrontSessions()

	cbFunc := func(s IServerSession) {
		log.Println("cb")
	}
	mgr.AddCloseCallback(cbFunc)
	mgr.AddCloseCallback(cbFunc)
	mgr.doCloseCallback(nil)
	
	id := mgr.findCloseCallback(cbFunc)
	log.Printf("find:%v\n", id)
	mgr.DelCloseCallback(id)
	mgr.doCloseCallback(nil)

	id = mgr.findCloseCallback(cbFunc)
	log.Printf("find:%v\n", id)
	mgr.DelCloseCallback(id)
	

	
}


