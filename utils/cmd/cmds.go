package cmd

import (
	"fmt"	
)

var (
	// 方便使用
	TheCmds *CmdMgr
)

func init() {
	fmt.Println("cmds.init")
	tryCreate()
}

func Visit() {
	tryCreate()
}

func tryCreate() {
	if TheCmds != nil {
		return
	}
	TheCmds = NewCmdMgr()
}

func RegisterCmd(c ICmd) {
	TheCmds.Register(c)
}

func DispatchCmd(content string) {
	TheCmds.Dispatch(content)
}
