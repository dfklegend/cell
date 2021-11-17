package cmd

import (
	"fmt"
	"strings"
)

type ICmd interface {
	GetName() string
	Do(args []string)
}

type CmdMgr struct {
	cmds map[string]ICmd
}

func NewCmdMgr() *CmdMgr {
	return &CmdMgr{cmds: make(map[string]ICmd)}
}

func (s *CmdMgr) Register(c ICmd) {
	fmt.Printf("register cmd:%s\n", c.GetName())
	s.cmds[(c).GetName()] = c
}

func (s *CmdMgr) Dispatch(content string) {
	subs := strings.Split(content, " ")
	if len(subs) < 1 {
		return
	}
	name := subs[0]
	c, ok := s.cmds[name]
	if !ok {
		return
	}
	c.Do(subs[1:])
}
