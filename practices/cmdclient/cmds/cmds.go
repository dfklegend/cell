package cmds

import (
	"fmt"

	cmd "github.com/dfklegend/cell/utils/libs/cmd"
)

var (
	TheCmds *cmd.CmdMgr
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
	TheCmds = cmd.NewCmdMgr()
}

func RegisterCmd(c cmd.ICmd) {
	TheCmds.Register(c)
}

func DispatchCmd(content string) {
	TheCmds.Dispatch(content)
}
