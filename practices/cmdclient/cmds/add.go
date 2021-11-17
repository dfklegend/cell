package cmds

import (
	"fmt"
	"strconv"
)

func init() {
	Visit()
	RegisterCmd(&AddCmd{})
}

type AddCmd struct {
}

func (s *AddCmd) GetName() string {
	return "add"
}

func (s *AddCmd) Do(args []string) {
	//fmt.Println("cmd1.Do")
	if len(args) < 2 {
		fmt.Println("need 2 args")
		return
	}
	a, _ := strconv.Atoi(args[0])
	b, _ := strconv.Atoi(args[1])
	c := a + b

	fmt.Printf("%v+%v=%v\n", a, b, c)
}
