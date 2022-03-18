package znet

import (
	"ZinxV0.2/src/zinx/utils"
	"ZinxV0.2/src/zinx/ziface"
	"fmt"
	"net"
)

// Server 定义IServer.go 的接口实现,定义一个server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器绑定的ip 版本号
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int

	//当前server的消息管理模块,用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandle

	//路由功能,给当前的服务注册一个路由方法,供客户端的链接处理使用
	//Router ziface.IRouter
}

// Start 启动服务器
func (s *Server) Start() {
	fmt.Println("配置信息", utils.GlobalObject)
	fmt.Printf("[start] IP:%s , Port:%d ,is starting ", s.IP, s.Port)

	go func() {
		//1.获取一个tcp的addr
		// ResolveTCPAddr返回TCP端点的地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error:", err)
			return
		}

		//2.监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, "err :", err)
			return
		}

		fmt.Println("start Zinx server succ,", s.Name, "succ,Listenning...")
		var cid uint32
		cid = 0

		//3.阻塞的等待客户端连接,处理客户端链接业务(读写)
		for true {
			//如果有连接会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//已经与客户端建立连接

			//将处理新连接的业务方法 和 conn进行绑定 得到链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()

}

/*//定义当前客户端链接所绑定handle api (目前这个handle是写死的,以后优化应该由用户指定此方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显
	_, err := conn.Write(data[:cnt])
	fmt.Println("CallBackToClient里 返回的是: ", string(data))
	if err != nil {
		fmt.Println("回显失败:", err)
		return errors.New("CallBackToClient err")
	}
	return nil
}
*/

// Stop 停止服务器
func (s *Server) Stop() {
	//todo 将一些服务器的资源,状态或者一些已经开辟的链路信息,进行停止或者回收
}

// Serve 运行服务器
func (s *Server) Serve() {
	//启动服务器
	s.Start()

	//TODO 做额外业务

	//阻塞状态
	select {}
}

//添加一个路由
func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)

}

// NewServer 初始化Server模块的方法
func NewServer(name string) ziface.IServer {
	server := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMagHandle(),
	}
	return server

}
