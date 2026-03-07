package siface

import "net"

//定义连接接口
type IConnection interface {
	Start()
	Stop()
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection()                 *net.TCPConn
	GetConnID()                        uint32
	RemoteAddr()                       net.Addr
	//直接将Message数据发送数据给远程的TCP客户端
	SendMsg(msgId uint32, data []byte) error
}