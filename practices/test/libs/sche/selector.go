package sche

import (	
	"fmt"
	"reflect"
	"time"
	"sync"
)

type ISelectorContainor interface {	
	Dummy()
}

type IChanSelector interface {
	GetChannel() reflect.Value
	DoTask(c ISelectorContainor, v reflect.Value, recvOk bool)
}

// -----------------------------
type SelectorData struct {
	selector IChanSelector	
	open bool;
}

// 动态的加入通道读取Select
// 避免写一个很明确的select列表
type MultiSelector struct {
	selectors []*SelectorData

	// 当前cases
	cases []reflect.SelectCase
	// 运行中的
	runnings []*SelectorData
	// 
	// 是否脏
	dirty bool

	mutex sync.Mutex
}

func NewMultiSelector() *MultiSelector {
	return &MultiSelector {
		selectors: make([]*SelectorData, 0),
		cases: make([]reflect.SelectCase, 0),
		runnings: make([]*SelectorData, 0),
		dirty: false,
	}
}

func (s *MultiSelector) Dummy() {
}

// 加入一个锁
func (s *MultiSelector) AddSelector(selector IChanSelector) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// TODO:判断是否已经加入过
	one := &SelectorData{
		selector: selector,
		open: true,
	}
	s.selectors = append(s.selectors, one)
	fmt.Printf("add selector %v\n", selector)
	fmt.Printf("selectors %v\n", s.selectors)

	s.markDirty();	
}

func (s *MultiSelector) markDirty() {
	s.dirty = true;
}

func (s *MultiSelector) makeCases() {
	fmt.Printf("makeCases \n")

	var cases = make([]reflect.SelectCase, 0)
	var runnings = make([]*SelectorData, 0);
	for i := 0; i < len(s.selectors); {
		one := reflect.SelectCase{}
		data := s.selectors[i]
		if(!data.open) {
			// TODO: 
			i ++
			continue;
		}

		//fmt.Printf("selector %v\n", selector)
		one.Dir = reflect.SelectRecv
		one.Chan = data.selector.GetChannel()

		cases = append(cases, one)
		runnings = append( runnings, data)

		i ++
	}

	// 运行中
	s.cases = cases;
	s.runnings = runnings;
}

func (s *MultiSelector) HandleOnce() {
	// 组织selectCases
    // 获取chosen
    // 获取chosen的runner
	// 调用runner.DoTask(c, v)		
	
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.dirty {
		s.makeCases()
		s.dirty = false
	}
	
	cases := s.cases
	// 调用reflect.Select size为0会一直卡住
	if(len(cases) == 0) {
		time.Sleep(10 * time.Millisecond)
		return;
	}	
	
	// 如果某个channel被关闭了,还在cases里会一直select到事件
	// 下次不能再加入到cases
	chosen, recv, recvOk := reflect.Select(cases)
	data := s.getSelectorByIndex(chosen);
	if(data == nil) {
		//
	}

	fmt.Printf("chosen:%v %v recv:%v recvOk:%v\n", chosen, data, recv, recvOk)
	data.selector.DoTask(s, recv, recvOk);

	// 处理被关闭的通道
	// 是否直接删除?
	if(!recvOk) {
		data.open = false
		s.markDirty()
	}
}

func (s *MultiSelector) getSelectorByIndex(tarIndex int) *SelectorData {
	return s.runnings[tarIndex]
}

// for test
func (s *MultiSelector) HandleLoop() {
	for true {
		s.HandleOnce()
	}
}

// -----------------------------
// 函数选择器
type SelectorFunc func(v reflect.Value, recvOk bool)

type FuncSelector struct {	
	fun SelectorFunc
	chanWait reflect.Value
}

func NewFuncSelector(chanWait reflect.Value, fun SelectorFunc) *FuncSelector {
	return &FuncSelector {
		fun: fun,
		chanWait: chanWait,
	}
}

func (s *FuncSelector) GetChannel() reflect.Value {
	return s.chanWait;
}

func (s *FuncSelector) DoTask(c ISelectorContainor, v reflect.Value, recvOk bool) {
	if(s.fun != nil) {
		s.fun(v, recvOk)
	}
}

