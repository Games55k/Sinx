package snet

import (
	"fmt"
	"net"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/sutils"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	msgHandler  siface.IMsgHandle
	ConnMgr     siface.IConnManager
	onConnStart func(conn siface.IConn)
	onConnStop  func(conn siface.IConn)
}

func NewServer() siface.IServer {

	sutils.GlobalObject.Reload()

	s := &Server{
		Name:       sutils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         sutils.GlobalObject.Host,
		Port:       sutils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
		onConnStart: func(conn siface.IConn) {},
		onConnStop:  func(conn siface.IConn) {},
	}
	return s
}

func (s *Server) Start() {
	fmt.Println("[Sinx] Server Name:", s.Name, "listenner at IP:", s.IP, " Port:", s.Port)

	fmt.Printf("[Sinx] Version: %s, MaxConn: %d,  MaxPacketSize: %d\n",
		sutils.GlobalObject.Version,
		sutils.GlobalObject.MaxConn,
		sutils.GlobalObject.MaxPacketSize)

	s.msgHandler.StartWorkerPool()

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("[Sinx] resolve tcp address err: ", err)
			return
		}

		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("[Sinx] listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("[Sinx] Listenning...")

		var cid uint32 = 0

		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("[Sinx] Accept err ", err)
				continue
			}

			if s.ConnMgr.Len() >= sutils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConntion(s, conn, cid, s.msgHandler)
			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[Sinx] stop server , name ", s.Name)

	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(msgId uint32, router siface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() siface.IConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(siface.IConn)) {
	s.onConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(siface.IConn)) {
	s.onConnStart = hookFunc
}

func (s *Server) GetOnConnStart() func(siface.IConn) {
	return s.onConnStart
}

func (s *Server) GetOnConnStop() func(siface.IConn) {
	return s.onConnStop
}

func (s *Server) GetMsgHandler() siface.IMsgHandle {
	return s.msgHandler
}