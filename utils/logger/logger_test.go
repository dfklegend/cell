package logger

import (	
	"testing"	

	"github.com/sirupsen/logrus"
)

func Test_Normal(t *testing.T) {
	SetDebugLevel()
	Log.Warn("warn")
	Log.Error("error")
	Log.Debug("debug1")
	SetInfoLevel()
	Log.Debug("debug2")
	Log.Error("error")
	SetLogLevel(logrus.DebugLevel)
	Log.Debug("debug3")
}


func FuncTest_Rotation(t *testing.T) {
	for i := 0; i < 10000000; i ++ {
		Log.Errorf("Error test\n")
	}
}