package common

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// 获取当前routine的id
// 用于测试
func GetRoutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	//str := string(buf[:n])
	//fmt.Println("str: %v", str)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
