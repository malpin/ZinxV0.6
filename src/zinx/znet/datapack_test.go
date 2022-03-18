package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//只是负责测试datapack拆包 封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/

	//1.创建sockerTCP
	listenner, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建一个go承载 负责从客户端处理业务
	go func() {
		//2.从客户端读取数据,拆包处理
		for true {
			conn, err := listenner.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}

			go func(conn net.Conn) {
				//处理客户端请求
				//拆包 读两次
				//定义一个拆包对象dp
				dp := NewDataPack()
				for true {
					// 1.读出包中的head
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println(err)
						return
					}
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println(err)
						return
					}

					//msg是有数据进行第二次读取
					if msgHead.GetMsgLen() > 0 {
						// 2.第二次从conn读,根据head中datalen 再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println(err)
							return
						}

						fmt.Println("读完消息: ")
						msg.ToString()
					}

				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
		return
	}
	//创建一个封包对象 dp
	dp := NewDataPack()
	//模拟粘包过程,封装两个msg一同发送
	//封装第一个msg
	message1 := &Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte{'q', 'w', 'q', 'w', 'q'},
	}
	sendDatal1, err := dp.Pack(message1)
	if err != nil {
		fmt.Println(err)
		return
	}
	//封装第二个msg
	message2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'z', 'c', 'v', 'r', 't'},
	}
	sendDatal2, err := dp.Pack(message2)
	if err != nil {
		fmt.Println(err)
		return
	}

	//将两个包粘在一起
	bytes := append(sendDatal1, sendDatal2...)

	//一次性发送给服务器
	_, err = conn.Write(bytes)
	if err != nil {
		return
	}

	//阻塞
	select {}
}
