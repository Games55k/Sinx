package snet

import (
	"fmt"
	"net"

	"github.com/Games55k/Sinx/siface"
)

// 客户端接口实现
type Client struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
	conn      siface.IConn
	// connMux     sync.Mutex
	onConnStart func(conn siface.IConn)
	onConnStop  func(conn siface.IConn)
	msgHandler  siface.IMsgHandle

	exitChan chan struct{}
}

var _ siface.IClient = (*Client)(nil)

func NewClient(name, ipVersion, ip string, port int) siface.IClient {
	c := &Client{
		Name:      name,
		IPVersion: ipVersion,
		IP:        ip,
		Port:      port,

		msgHandler: NewMsgHandle(),
	}
	return c
}

func (c *Client) Restart() {
	c.Stop()
	c.Start()
}

func (c *Client) Start() {
	c.exitChan = make(chan struct{})
	go c.msgHandler.StartWorkerPool()

	go func() {
		addr, err := net.ResolveTCPAddr(c.IPVersion, fmt.Sprintf("%s:%d", c.IP, c.Port))
		if err != nil {
			fmt.Println("[Cinx] resolve tcp address err: ", err)
			return
		}

		conn, err := net.DialTCP(c.IPVersion, nil, addr)
		if err != nil {
			fmt.Println("[Cinx] dial tcp err: ", err)
		}

		c.conn = NewClientConn(c, conn)

		go c.conn.Start()

		select {
		case <-c.exitChan:
			fmt.Println("[Sinx] client exit")
			return
		}
	}()
}
func (c *Client) Stop() {
	con := c.Conn()
	if con != nil {
		con.Stop()
	}
}
func (c *Client) Conn() siface.IConn {
	return c.conn
}
func (c *Client) AddRouter(msgId uint32, router siface.IRouter) {
	c.msgHandler.AddRouter(msgId, router)
}
func (c *Client) GetMsgHandler() siface.IMsgHandle {
	return c.msgHandler
}
func (c *Client) SetOnConnStart(f func(siface.IConn)) {
	c.onConnStart = f
}
func (c *Client) SetOnConnStop(f func(siface.IConn)) {
	c.onConnStop = f
}
func (c *Client) GetOnConnStart() func(siface.IConn) {
	return c.onConnStart
}
func (c *Client) GetOnConnStop() func(siface.IConn) {
	return c.onConnStop
}