# 兼容性说明
	由于etcd，gRPC版本要求 v1.26.0

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

	
	