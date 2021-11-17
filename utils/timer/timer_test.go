package timer

import (
	"log"
	"testing"
	"time"
	//"sync/atomic"

	"github.com/stretchr/testify/assert"

	"github.com/dfklegend/cell/utils/common"
)

// 测试timer基本流程
// 创建,cancel
func Test_Normal(t *testing.T) {
	mgr := GetTimerMgr()
	
	total := 0
	t1R := 0
	t2R := 0

	t2 := mgr.AddTimer(100*time.Millisecond, func(args ...interface{}) {
		t1R = common.GetRoutineID()
		log.Printf("t2 r:%v %v\n", t1R, args)
		
		total += 2
	}, 3, 5)


	t1 := mgr.After(350*time.Millisecond, func(args ...interface{}) {
		t2R = common.GetRoutineID()
		log.Printf("t1 r:%v\n", t2R)		
		total += 1
		mgr.Cancel(t2)
	}, 3)

	run := true
	go func() {
		for run {
			select {
			case t := <-mgr.GetQueue():
				mgr.Do(t)
			}
		}
	}()

	time.Sleep(1 * time.Second)
	run = false
	mgr.Cancel(t1)		

	assert.Equal(t, 7, total, "not equal")
	assert.Equal(t, t1R, t2R, "same routine")
}
