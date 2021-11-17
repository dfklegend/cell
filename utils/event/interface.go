package event


/**
 * 提供两级的事件中心
 * Global事件，可以推送到其他local
 * Local可以注册局部以及全局事件
 * Local事件，需使用者保证线程安全
 * (也可以使用SetLocalUseChan来确保局部事件也使用channel来
 * 拉取事件(能定制运行coroutine))
 * 并且需要自主拉取global事件 select ChanEvent
 */

type ChanEvent chan *EventObj
type CBFunc func(args ...interface{})

type EventObj struct {
	EventName 	string
	Args 		[]interface{}
}

type EventListener struct {
	Id 			uint64
	// 自带参数	
	Args 		[]interface{}
	CB 			CBFunc
}

// 全局的事件中心
type IGlobalEventCenter interface {
	Subscribe(eventName string, child ILocalEventCenter)
	Unsubscribe(eventName string, child ILocalEventCenter)
	Publish(eventName string, args ...interface{})	
}

// 本地事件中心，为本地服务，用chan拉取到保证cb的执行体
type ILocalEventCenter interface {
	// 自身id
	GetId() uint64
	Clear()
	// 本地事件
	Subscribe(eventName string, cb CBFunc, args ...interface{}) uint64
	Unsubscribe(eventName string, id uint64)
	// 全局事件
	GSubscribe(eventName string, cb CBFunc, args ...interface{}) uint64
	GUnsubscribe(eventName string, id uint64)

	GetChanEvent() ChanEvent
	DoEvent(e *EventObj)

	Publish(eventName string, args ...interface{})
}
