package main


import (
	"fmt"
	"encoding/json"
)

func Test1() {
    m := make(map[string]interface{})

    m["v1"] = "str"
    m["v2"] = 1

    fmt.Printf("%+v\n", m)
    r, _ := json.Marshal(m)
    result := string(r)
    fmt.Println(result)

    var nm interface{}
    json.Unmarshal(r, &nm)
    fmt.Printf("nm:%+v\n", nm)
}

// 测试反序列化
// 存在的key，覆盖(包括类型)
// 不存在的key,添加
// 字符串没包含的key,不受影响
func Test2() {
    m := make(map[string]interface{})

    m["v1"] = "str"
    m["v2"] = 1
    m["v3"] = 1

    fmt.Printf("%+v\n", m)
    m1 := make(map[string]interface{})
    
    m1["v2"] = 3
    m1["v3"] = "str"
    m1["v4"] = "str"

    r, _ := json.Marshal(m1)
    
    json.Unmarshal(r, &m)
    fmt.Printf("m:%+v\n", m)
}

// 测试数组
type Exam struct {
    Array []int `json:"array"`
}

func NewExam() *Exam {
    return &Exam {
        Array: make([]int, 0), 
    }
}

func Test3() {
    e := NewExam()
    e.Array = append(e.Array, 0, 1, 2)

    s, _ := json.Marshal(e)
    fmt.Printf("%+v\n", string(s))

    e1 := NewExam()
    json.Unmarshal(s, &e1)

    fmt.Printf("%+v\n", e1)
}

func main() {
	Test3()
}
