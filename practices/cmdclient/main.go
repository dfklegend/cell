package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"

	"net/http"
	_ "net/http/pprof"

	"github.com/dfklegend/cell/cmdclient/cmds"
	"github.com/dfklegend/cell/utils/libs/cmd"
	"github.com/dfklegend/cell/utils/libs/util"
)

func main() {
	fmt.Println("cmdclient start")
	fmt.Println(util.GetRoutineID())

	testI()
	cmds.Visit()

	go processInput()
	go webServe()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	s := <-c
	fmt.Println("Got signal:", s)
}

func webServe() {
	http.ListenAndServe("0.0.0.0:6060", nil)
}

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
		fmt.Print("input:")
		s := getInput()
		fmt.Println(s)

		// 传到cmdMgr中
		cmds.DispatchCmd(s)
	}
}

func testI() {
	c1 := cmds.AddCmd{}
	var i cmd.ICmd
	i = &c1

	fmt.Printf("%p %p %p %p\n", c1, &c1, &i, i)
}
