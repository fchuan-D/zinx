package znet

import (
	"fmt"
	"net"
	"zinx/src/zinx/utils"
	"zinx/src/zinx/ziface"
)

//IServer 的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器名称
	Name string
	//服务器监听IP版本
	IPVersion string
	//服务器监听IP
	IP string
	//服务器监听端口Port
	Port int
	// 当前Server注册的Router
	Router ziface.IRouter
}

//启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listenner at IP:%s Port:%d is starting \n", s.IP, s.Port)

	go func() {
		//获取一个TCP的Address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve TCP addr err:", err)
			return
		}

		//监听服务器的地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "	err:", err)
			return
		}

		fmt.Println("start Zinx server succ,", s.Name, "succ,Listenning...")
		var cid uint32 = 0

		//阻塞的等待客户端连接，处理客户端请求
		for {
			//如果有客户端请求连接，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//将server注册的router封装到新的连接对象
			//得到封装以后的Conn连接对象
			pakConn := NewConnection(conn, cid, s.Router)
			cid++
			//启动当前连接的业务处理
			go pakConn.Start()
		}
	}()
}

//运行服务器
func (s *Server) Serve() {
	//启动Serve的服务功能
	s.Start()

	//TODO:启动服务后的额外业务拓展

	//阻塞状态
	select {}
}

//关闭服务器
func (s *Server) Stop() {
	//TODO:将一些连接的资源或信息回收

}

//路由功能，给当前Server添加一个Router
func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router succ!!")
}

/*
	初始化Server模块
*/
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		Router:    nil,
	}
	return s
}
