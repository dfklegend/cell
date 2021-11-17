package main

import (
	"fmt"
	"net"
	"time"
	"encoding/binary"
)

func makePacket(data []byte) []byte {
	var msgLen uint32
	msgLen = uint32(len(data))

	msg := make([]byte, 4+msgLen)
	binary.LittleEndian.PutUint32(msg, msgLen)
	
	copy(msg[4:], data)	
	return msg;
}

func writePacket(conn net.Conn, data []byte) {
	msg := makePacket(data)
	conn.Write(msg)
}

// 发送2个字节(长度的一半)
func test1(conn net.Conn, data []byte) {
	msg := makePacket(data)

	// 估计先发2字节
	conn.Write(msg[:2])
	time.Sleep(8 * time.Second)
	conn.Write(msg[2:])
	//
}

func main() {
    conn, err := net.Dial("tcp", "127.0.0.1:12345")
    if err != nil {
        fmt.Println("连接服务器失败, err", err)
        return
    }
    defer conn.Close()
    for i := 0; i < 20; i++ {
        msg := `Hello, Hello. How are you?`
        //conn.Write([]byte(msg))
        //writePacket(conn, []byte(msg))
        test1(conn, []byte(msg))
    }
}