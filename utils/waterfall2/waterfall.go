package waterfall

import (
	scheP "github.com/dfklegend/cell/utils/sche"
)

// 回调函数
type Callback func(err bool, args ...interface{})
type Task func(args ...interface{})

// 实现一个类似nodejs waterfall
// 每个回调都在调用routine中调用

type Chain struct {
	node *WFNode
	// 任务列表
	tasks []Task
	// 最后的任务
	final   Task
	cursor  int
	curArgs []interface{}
	// cb函数
	cb           Callback
	startTime    int
	totalTimeout int
}

func (s *Chain) _exec(index int, args ...interface{}) {
	args = append([]interface{}{s.cb}, args...)
	s.tasks[index](args...)
}

func (s *Chain) tryExec(index int, args ...interface{}) {
	if index < len(s.tasks) {
		s._exec(index, args...)
	} else {
		args = append([]interface{}{false}, args...)
		s.doFinal(args...)
	}
}

func (s *Chain) doNext(args ...interface{}) {
	s.cursor++
	s.tryExec(s.cursor, args...)
}

func (s *Chain) doFinal(args ...interface{}) {
	s.final(args...)
}

func (s *Chain) doCallback(err bool, args ...interface{}) {
	if err {
		args = append([]interface{}{true}, args...)
		s.doFinal(args...)
		return
	}

	// doNext
	s.doNext(args...)
}

// 代表
type WFNode struct {
	sche *scheP.Sche
	//chains   []Chain
	//chanNext chan *Chain
}

func Waterfall(sche *scheP.Sche, tasks []Task, final Task) *Chain {
	thisChain := &Chain{
		tasks: tasks,
		final: final,
	}

	thisChain.cb = func(err bool, args ...interface{}) {
		sche.Post(func() (interface{}, error) {
			thisChain.doCallback(err, args...)
			return 0, nil
		})
	}

	sche.Post(func() (interface{}, error) {
		thisChain.tryExec(0)
		return 0, nil
	})
	return thisChain
}
