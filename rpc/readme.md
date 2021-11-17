# apientry
    collection
        注册接口注册
        统一访问        
        arg和cbFunc里返回的outArg都是流数据
        Call(route string, args []byte, cbFunc HandlerCBFunc, ext interface{}) error

    container
        接口的容器，提供具体handler处理
        负责将参数数据流反序列化成参数
        并将返回值序列化成数据流