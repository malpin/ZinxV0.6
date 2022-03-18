package main

import (
	_ "ZinxV0.2/src/zinx/utils"
	"ZinxV0.2/src/zinx/ziface"
	"ZinxV0.2/src/zinx/znet"
	"fmt"
)

//自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//处理conn业务的主方法hook
func (ping *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端的数据 再回写
	fmt.Println("server读到了 ZinxV0.5 request msgid= ", request.GetMsgID(),
		"msgLen= ", request.GetMsgLen(),
		"data= ", request.GetData())

	err := request.GetConnection().SendMsg(200, []byte("ping..haha"))
	fmt.Println("发了")
	if err != nil {
		fmt.Println(err)
	}

}

//自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

//处理conn业务的主方法hook
func (ping *HelloRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call  HelloRouter  Handle...")
	//先读取客户端的数据 再回写
	fmt.Println("server读到了 ZinxV0.5 request msgid= ", request.GetMsgID(),
		"msgLen= ", request.GetMsgLen(),
		"data= ", request.GetData())

	err := request.GetConnection().SendMsg(400, []byte("hello..haha"))
	fmt.Println("发了")
	if err != nil {
		fmt.Println(err)
	}

}

func main() {
	//创建一个server句柄 使用zinx的api
	s := znet.NewServer("[zinx v0.6]")
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	//启动server
	s.Serve()
}
