package debug

import (
	"log"
	"runtime"
)

func DumpCurInfo() {
	log.Printf("curRoutineNum:%v\n", runtime.NumGoroutine())
}
