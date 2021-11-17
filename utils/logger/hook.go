package logger

import (
    "time"
    "fmt"

    rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"	
    "github.com/rifflock/lfshook"
)

func InitHook(Log *logrus.Logger) {
    pathMap := lfshook.PathMap{
        logrus.InfoLevel:  "./info.log",
        logrus.ErrorLevel: "./error.log",
    }

    Log.Hooks.Add(lfshook.NewHook(
        pathMap,
        &logrus.JSONFormatter{},
    ))
}


func InitRotationHook(Log *logrus.Logger) {
    path := "./go.%v.log"
    writer, _:= rotatelogs.New(
        fmt.Sprintf(path,"%Y%m%d%H%M"),
        rotatelogs.WithLinkName(path),
        rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
        rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
    )

    writerMap := lfshook.WriterMap{
        logrus.InfoLevel:  writer,
        logrus.ErrorLevel: writer,
    }

    Log.Hooks.Add(lfshook.NewHook(
        writerMap,
        &logrus.JSONFormatter{},
    ))
}