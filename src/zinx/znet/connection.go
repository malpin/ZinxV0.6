package znet

import (
	"ZinxV0.2/src/zinx/ziface"
	"errors"
	"fmt"
	"io"
	"net"
)

//链接模块
type Connection struct {
	//当前链接的 socket TCP 套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//当前链接所绑定的处理业务方法API
	//handleAPI ziface.HandleFunc

	//告知当前链接已经 退出/停止 的 channel
	ExitChan chan bool

	//无缓冲通道,用于读,写 Goroutine 之间的消息通信
	msgChan chan []byte

	//MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//该链接处理的方法Router
	//Router ziface.IRouter
}

//初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}
	return c
}

//写消息Goroutine ,专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[Writer exit!]")

	//阻塞等待channel的消息 进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("sand data error", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经退出 此时Writer也要退出
			return
		}
	}

}

//启动链接 让当前链接准备工作
func (c *Connection) Start() {
	fmt.Println("链接模块 Conn Start() ConnID=  ", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//启动从当前链接写数据的业务
	go c.StartWriter()
}

//停止链接 结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop() ConnID=  ", c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket链接
	c.Conn.Close()
	//告知Writer关闭
	c.ExitChan <- true

	//关闭管道
	close(c.ExitChan)
	close(c.msgChan)
}

//获取当前链接的绑定 套接字socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前链接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//获取远程客户端的TCP状态IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//发送数据,将数据发送给远程的客户端  此方法 将发送给客户端的数据先封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection close when send msg")
	}

	//将data进行封包
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data)) //binaryMsg 为二进制
	if err != nil {
		fmt.Println("Pack err", err)
		return errors.New("pack err")
	}

	//将数据发送给 写消息Goroutine ,专门发送给客户端消息的模块
	c.msgChan <- binaryMsg

	/*if _, err := c.Conn.Write(binaryMsg); err != nil {
		fmt.Println("Write msg id:", msgId, " err", err)
		return errors.New("Conn Write err ")
	}*/
	return nil
}

//读数据
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit,remote addr is(读取器退出，远程地址为)", c.RemoteAddr().String())
	defer c.Stop()

	//读取客户端的数据到buf 中 最大512字节
	for {
		//v0.5注释掉了
		/*buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		fmt.Printf("链接模块 读到了%s \n", string(buf))
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}*/

		//创建一个拆包解包的对象
		dp := NewDataPack()
		//读取客户端的msg head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("读取MSG头错误 read msg head error", err)
			break
		}
		//拆包 得到msgID msgdatalen 放入msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("Unpack err", err)
			break
		}
		//根据datalen再次读取 放在 msg.data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err", err)
				break
			}
		}
		msg.SetData(data) //读出的数据放入msg.data

		//得到当前conn数据和Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//从路由中找到注册绑定的conn对应router调用
		//根据绑定好的MsgID 找到对应处理api业务 执行
		go c.MsgHandler.DoMsgHandler(&req)

		//从路由中得到注册绑定的Conn对应的router调用
		//go func(request ziface.IRequest) {
		//	c.Router.PreHandle(request)
		//	c.Router.Handle(request)
		//	c.Router.PostHandle(request)
		//}(&req)

		//调用当前链接所绑定的HandleAPI
		//err = c.handleAPI(c.Conn, buf, cnt)
		//if err != nil {
		//	fmt.Println("ConnID ", c.ConnID, " handle is error", err)
		//	break
		//}
	}
}
