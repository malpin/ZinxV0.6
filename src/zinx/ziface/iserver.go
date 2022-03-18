package ziface

// IServer 定义服务器接口
type IServer interface {

	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Serve 运行服务器
	Serve()

	//路由功能,给当前的服务注册一个路由方法,供客户端的链接处理使用
	AddRouter(uint32, IRouter)
}
