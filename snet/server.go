package snet

import (
	"fmt"
	"net"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/sutils"
)

// iServer 接口实现，定义一个Server服务类
type Server struct {
	Name       string
	IPVersion  string
	IP         string
	Port       int
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler siface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr siface.IConnManager
	//该Server的连接创建时Hook函数
	OnConnStart func(conn siface.IConn)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn siface.IConn)
}

// 开启 server 服务（无阻塞）
func (s *Server) Start() {
	fmt.Println("[Sinx] Server Name:", s.Name, "listenner at IP:", s.IP, " Port:", s.Port)

	fmt.Printf("[Sinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		sutils.GlobalObject.Version,
		sutils.GlobalObject.MaxConn,
		sutils.GlobalObject.MaxPacketSize)

	// 创建协程不间断处理链接
	go func() {
		//开启工作池
		s.msgHandler.StartWorkerPool()
		//封装 tcp 地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[Sinx] resolve tcp address err: ", err)
			return
		}

		//创建监听 socket
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[Sinx] listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("[Sinx] start success, now listenning...")

		//简单实现一个自增的连接 ID
		var cid uint32 = 0

		//持续监听客户端连接
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("[Sinx] Accept err ", err)
				continue
			}

			//3.2 判断当前服务器的连接数是否已经超过最大连接数
			if s.ConnMgr.Len() >= sutils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			//3.3 初始化连接模块
			dealConn := NewConntion(s, conn, cid, s.msgHandler)
			cid++

			//3.4 启动协程处理当前连接的业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Sinx] stop server , name ", s.Name)

	//通过 ConnManager 清除并停止所有连接
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢

	//阻塞
	select {}
}

// 为特定消息注册处理函数
func (s *Server) AddRouter(msgId uint32, router siface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)

	fmt.Println("[Sinx] Add Router success! ")
}

func (s *Server) GetConnMgr() siface.IConnManager {
	return s.ConnMgr
}

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(siface.IConn)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(siface.IConn)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn siface.IConn) {
	if s.OnConnStart != nil {
		fmt.Println("[Sinx] CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn siface.IConn) {
	if s.OnConnStop != nil {
		fmt.Println("[Sinx] CallOnConnStop....")
		s.OnConnStop(conn)
	}
}

func NewServer() siface.IServer {
	//初始化全局配置文件
	sutils.GlobalObject.Reload()
	s := &Server{
		Name:       sutils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         sutils.GlobalObject.Host,
		Port:       sutils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}

func (s *Server) GetMsgHandler() siface.IMsgHandle {
	return s.msgHandler
}