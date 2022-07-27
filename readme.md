# 整体目标
    使用go搭建服务器，希望满足以下需求
    . 架构参考pomelo
    . 能方便的切分服务器功能，组织子服务器
    . 易于使用的服务器rpc
    . 服务发现
    . 与逻辑开发无关
    . web控制台
    . 进程热替换机制

    go语言相关
    . 定义逻辑的routine调度原则
    . 提供对应的支持库

# 整体进展
## 框架部分
	. RPC基础功能
		gRPC或者TCP协议
		Json格式
	. 服务器框架功能
		自由定制服务器组件
	. 前端协议路由转发
	. 前端同时允许TCP和WS接入
	. 基于ETCD的服务发现
	. 整理了调度器机制，可以方便的将
	  任务推送到指定调度器执行
	. 整理了通用的Timer机制
		携程和调度器驱动
	. 整理了通用的Event机制
		调度器内部事件
		全局事件

## 实践部分
	. 聊天室和roll点功能
	. 尝试搭建MMO结构，进行一些测试

# TODO
	. 更好的性能分析
	. 进一步搭建模拟MMO工程，验证合适的模型




# server/examples
    需要安装etcd
    https://github.com/etcd-io/etcd/releases
    https://github.com/etcd-io/etcd/releases/download/v3.4.19/etcd-v3.4.19-windows-amd64.zip
## chat
### 功能
    可以通过浏览器启动chat-client/index.html来启动客户端
    客户端启动后，登录服务器，随机分配到某个聊天服务器
    并加入房间，房间最多3人，可以聊天
    聊天输入/roll可以掷骰子
    玩家进入服务器总是进入最前面的房间
    房间不会删除
### 配置与启动
    data/config/servers.yaml
    配置了多个服务器
    使用id来启动不同的服务器类型
    ./chat.exe -id=gate-1
    ./chat.exe -id=chat-1

    最少需要包含一个gate和一个chat组件

## 客户端
    chat-client
    网页客户端 index.html拖拽到浏览器执行

    chat-client-go
    go 客户端



# refs
https://github.com/topfreegames/pitaya
https://github.com/lonng/nano