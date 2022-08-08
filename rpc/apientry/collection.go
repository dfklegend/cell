package apientry

import (
	"fmt"
	"strings"

	"github.com/dfklegend/cell/utils/logger"
)

// 在启动时集中注册服务
// map就不存在并发问题
func NewCollection() *APICollection {
	collection := &APICollection{
		Containers:  make(map[string]*APIContainer),
		Entries:     &APIEntries{},
		apiFormater: defaultFormater,
	}
	return collection
}

type APICollection struct {
	Containers  map[string]*APIContainer
	Entries     *APIEntries
	apiFormater IAPIFormater
}

func (self *APICollection) SetFormater(formater IAPIFormater) {
	self.apiFormater = formater
}

func (self *APICollection) Register(e IAPIEntry, options ...Option) {
	self.Entries.Register(e, options...)
}

func (self *APICollection) Build() error {
	// 初始化
	//clear old
	self.Containers = make(map[string]*APIContainer)

	entries := self.Entries.List()
	for _, c := range entries {
		err := self.newService(c.Entry, c.Opts...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *APICollection) newService(comp IAPIEntry, opts ...Option) error {
	srv := NewContainer(comp, opts...)

	if _, ok := self.Containers[srv.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", srv.Name)
	}

	if err := srv.ExtractHandler(self.apiFormater); err != nil {
		return err
	}

	// 记录
	self.Containers[srv.Name] = srv
	logger.Log.Infof("newService:%v \n", srv.Name)
	return nil
}

// Call
/**
 * 1. 转发给对应的Service
 * 2. 按目标参数类型反序列化
 * cb(error, []byte result)
 * 出入口都是[]byte (json序列化)
 * ext
 * 		一个自定义参数，目前是前端接口的session
 *
 *
 * 将请求转发给底层的各接口执行
 * @param route{string} 服务.method路由字符串
 * @param args{[]byte} 请求参数流化的数据(json)
 * @param cbFunc{HandlerCBFunc} 执行完后的回调
 *        func(e error, outArg interface{})
 *        		error 是否执行错误
 *        		outArt json序列化的返回值
 * @param ext{interace{}} 自定义参数，server:前端接口是session
 */
func (self *APICollection) Call(route string, args []byte,
	cbFunc HandlerCBFunc, ext interface{}) error {
	subs := strings.Split(route, ".")
	if len(subs) != 2 {
		return fmt.Errorf("bat route: %s", route)
	}
	serviceName := subs[0]
	methodName := subs[1]

	service := self.Containers[serviceName]
	if service == nil {
		logger.Log.Warnf("can not find service: %s", serviceName)
		return fmt.Errorf("can not find service: %s", serviceName)
	}
	return service.CallMethod(self.apiFormater, methodName, args, cbFunc, ext)
}

func (self *APICollection) HasMethod(route string) bool {
	return true
}

// 序列化
func (s *APICollection) AddReqReply(reqId int, args interface{}) {

}
