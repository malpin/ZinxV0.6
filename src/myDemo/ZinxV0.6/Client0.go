package main

import (
	"ZinxV0.2/src/zinx/znet"
	"fmt"
	"io"
	"net"
	"time"
)

//模拟客户端
func main() {
	fmt.Println("client0 start..")
	//1.直接链接远程服务器,得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err ,exit ")
		return
	}

	//2.链接调用Write 写数据
	for true {
		//发送封包的message消息 msgId是0的消息 zinxV0.5 msg
		dp := znet.NewDataPack()
		binaryMag, err := dp.Pack(znet.NewMsgPackage(1, []byte("zinxV0.6 client0 msg")))
		if err != nil {
			fmt.Println(" dp.Pack err ", err)
			return
		}
		//发送
		_, err = conn.Write(binaryMag)
		if err != nil {
			fmt.Println("conn.Write err ", err)
			return
		}

		//服务器回复一个message数据    msgID=1, data=[]byte("ping..haha")
		//先读取流中的head部分,得到ID和datalen
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("ReadFull err", err)
			continue
		}

		//将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("Unpack err", err)
			continue
		}
		if msgHead.GetMsgLen() > 0 { //有数据
			//再根据datalen 进行二次读取data
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("io.ReadFull err", err)
				continue
			}
			fmt.Println("client读到了")
			msg.ToString()
		}

		//阻塞 防止cpu跑满
		time.Sleep(3 * time.Second)
	}
}
