package siface

import "net"

//定义连接接口
type IConnection interface {
	Start()
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn
	GetConnID()        uint32
	RemoteAddr()       net.Addr
}

//定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error