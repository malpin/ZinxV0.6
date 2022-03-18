package main

import "ZinxV0.2/src/zinx/znet"

func main() {
	//创建一个server句柄 使用zinx的api
	s := znet.NewServer("[zinx v0.2]")
	//启动server
	s.Serve()
}
