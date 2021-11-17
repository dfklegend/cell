package cmd

import (
	"os"
	"fmt"	
	"bufio"
)

func getInput() string {
	in := bufio.NewReader(os.Stdin)
	str, _, err := in.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(str)
}

func processInput() {
	for true {
		fmt.Print("please input:")
		s := getInput()		
		fmt.Println(s)
		// 传到cmdMgr中
		DispatchCmd(s)
	}
}

func StartConsoleCmd() {
	go processInput()
}