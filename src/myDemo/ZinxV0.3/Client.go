package main

import (
	"fmt"
	"net"
	"time"
)

//模拟客户端
func main() {
	fmt.Println("client start..")
	//1.直接链接远程服务器,得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err ,exit ")
		return
	}

	//2.链接调用Write 写数据
	for true {
		_, err := conn.Write([]byte("啊哈哈")) //字节数write  发送
		if err != nil {
			fmt.Println("Write conn err", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf) // cnt长度
		if err != nil {
			fmt.Println("Read conn err", err)
			return
		}
		fmt.Printf("客户端 接收到:%s,长度:%d \n", string(buf), cnt)

		//阻塞 防止cpu跑满
		time.Sleep(3 * time.Second)
	}
}
