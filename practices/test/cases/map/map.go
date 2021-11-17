package main


import (
	"fmt"
	"time"
)

// 测试map并发问题
// hashWriting
// 不写入(删除)，纯读取，可以并发


func GoLoop(f func(), num int) {
    for i := 0; i < num; i ++ {
        go func() {
            f()
        }()
    }
}

func main() {
    m := make(map[int]int)
    m[1] = 1

	GoLoop(func(){
        for true {
            //m[2] = 2
        }
    }, 20)

    GoLoop(func(){
        for true {
            _ = m[1]
        }
    }, 100)

    GoLoop(func(){
        for true {
            //delete(m, 2)
        }
    }, 20)
    time.Sleep(10*time.Second)
    fmt.Println("over")
}
