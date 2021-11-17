package interfaces

var Runtime *RuntimeData = &RuntimeData{}

// 运行时一些接口，避免交叉引用
type RuntimeData struct {    
    App IApp
}

func SetApp(app IApp) {
    Runtime.App = app
}

func GetApp() IApp {
    return Runtime.App
}