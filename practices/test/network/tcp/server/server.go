package server

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"

	MyIO "github.com/dfklegend/cell/practices/test/libs/io"
)

func packetRead(conn net.Conn) ([]byte, error) {
	var b [4]byte

	bufMsgLen := b[:4]
	// read len
	if _, err := MyIO.ReadFull(conn, bufMsgLen); err != nil {
		fmt.Println("读取客户数据失败，err:", err)
		return nil, err
	}

	msgLen := binary.LittleEndian.Uint32(bufMsgLen)
	// read后面的
	msgData := make([]byte, msgLen)

	if _, err := MyIO.ReadFull(conn, msgData); err != nil {
		fmt.Println("读取客户数据失败，err:", err)
		return nil, err
	}

	fmt.Println("收到client发来的数据：", string(msgData))

	return msgData, nil
}

func process2(conn net.Conn) {
	defer conn.Close()

	for {
		data, err := packetRead(conn)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取客户数据失败，err:", err)
			break
		}
		recvData := string(data)
		fmt.Println("收到client发来的数据：", recvData)
	}
}

func process1(conn net.Conn) {
	defer conn.Close()
	// 使用bufio的读缓冲区（防止系统缓冲区溢出）
	reader := bufio.NewReader(conn)
	var buf [1024]byte
	for {
		n, err := reader.Read(buf[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("读取客户数据失败，err:", err)
			break
		}
		recvData := string(buf[:n])
		fmt.Println("收到client发来的数据：", recvData)
	}
}

func process(conn net.Conn) {
	process2(conn)
}

func ListenAndServe() {
	listen, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		fmt.Println("监听失败, err:", err)
		return
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("建立会话失败, err:", err)
			continue
		}
		//conn.Set
		// 单独建一个goroutine来维护客户端的连接（不会阻塞主线程）
		go process(conn)
	}
}
