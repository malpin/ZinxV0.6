package main

import (
	"ZinxV0.2/src/zinx/ziface"
	"ZinxV0.2/src/zinx/znet"
	"fmt"
)

//自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//处理conn业务之前的钩子方法hook
func (ping *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping err")
	}
}

//处理conn业务的主方法hook
func (ping *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping... ping...\n"))
	if err != nil {
		fmt.Println("call back  ping err")
	}
}

//处理conn业务之后的钩子方法hook
func (ping *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back  after ping err")
	}
}

func main() {
	//创建一个server句柄 使用zinx的api
	s := znet.NewServer("[zinx v0.3]")
	s.AddRouter(&PingRouter{})
	//启动server
	s.Serve()
}
