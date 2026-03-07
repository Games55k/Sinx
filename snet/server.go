package snet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/sutils"
)

// iServer 接口实现，定义一个Server服务类
type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler siface.IMsgHandle
}

// ============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	fmt.Println("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

//============== 实现 ciface.IServer 里的全部接口方法 ========

func (s *Server) Start() {
	fmt.Printf("[START] Server listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	fmt.Printf("[Sinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		sutils.GlobalObject.Version,
		sutils.GlobalObject.MaxConn,
		sutils.GlobalObject.MaxPacketSize)

	go func() {
		//启动worker工作池机制
		s.msgHandler.StartWorkerPool()
		//获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		//2 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Sinx server  ", s.Name, " succ, now listenning...")

		//TODO server.go 应该有一个自动生成ID的方法

		var cid uint32
		cid = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}

			//3.2 TODO Server.Start() 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接

			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConntion(conn, cid, s.msgHandler)
			cid++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Cinx server , name ", s.Name)

	//TODO  Server.Stop() 将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()

	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId uint32, router siface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)

	fmt.Println("Add Router succ! ")
}

/*创建一个服务器句柄*/
func NewServer() siface.IServer {
	//先初始化全局配置文件
	sutils.GlobalObject.Reload()
	s := &Server{
		Name:       sutils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         sutils.GlobalObject.Host,
		Port:       sutils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
	}
	return s
}