package component

// 定制app的启动，关闭过程
type IComponent interface {
	Init()
	Start()
	Stop()
}

type DefaultComponent struct {
}

func (self *DefaultComponent) Init() {
}

func (self *DefaultComponent) Start() {
}

func (self *DefaultComponent) Stop() {
}
