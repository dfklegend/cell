package io

import (
	"fmt"
	sysIO "io"
	"runtime"
)

func ReadFull(r sysIO.Reader, buf []byte) (n int, err error) {
	return ReadAtLeast(r, buf, len(buf))
}

func ReadAtLeast(r sysIO.Reader, buf []byte, min int) (n int, err error) {
	if len(buf) < min {
		return 0, sysIO.ErrShortBuffer
	}

	nullTimes := 0
	for n < min && err == nil {
		var nn int

		// 结果Read没有数据时会阻塞routine,并不会多次
		// 空转
		fmt.Printf("read in\n")
		nn, err = r.Read(buf[n:])
		fmt.Printf("read out\n")
		n += nn
		fmt.Printf("nn:%d", nn)
		if nn == 0 && err == nil {
			nullTimes++
			fmt.Printf("nullTimes:%d", nullTimes)
			if nullTimes >= 10 {
				runtime.Gosched()
			}
		}
	}
	if n >= min {
		err = nil
	} else if n > 0 && err == sysIO.EOF {
		err = sysIO.ErrUnexpectedEOF
	}
	return
}
