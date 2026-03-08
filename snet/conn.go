package snet

import (
	"fmt"
	"net"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/sutils"
	
)

type Connection struct {
	TcpServer siface.IServer
	Conn      *net.TCPConn
	ConnID    uint32
	isClosed  bool

	IsClosed     atomic.Bool
	IsAborted    atomic.Bool
	IsClosedOnce atomic.Bool

	MsgHandler   siface.IMsgHandle
	ExitBuffChan chan struct{}

	msgChan      chan []byte
	msgBuffChan  chan []byte

	writerClosedChan chan struct{}

	property     map[string]interface{}
	propertyLock sync.RWMutex

	onConnStart func(conn siface.IConn)
	onConnStop  func(conn siface.IConn)

	wg          *sync.WaitGroup
}

func NewConntion(server siface.IServer, conn *net.TCPConn, connID uint32, msgHandler siface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan struct{}, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, sutils.GlobalObject.MaxMsgChanLen),
		writerClosedChan: make(chan struct{}),
		property:     make(map[string]interface{}),
		onConnStart:  server.GetOnConnStart(),
		onConnStop:   server.GetOnConnStop(),
		wg:           &sync.WaitGroup{},
	}
	c.IsClosed.Store(false)
	c.IsAborted.Store(false)
	c.IsClosedOnce.Store(false)

	server.GetConnMgr().Add(c)
	return c
}

func NewClientConn(client siface.IClient, conn *net.TCPConn) siface.IConn {
	c := &Connection{
		TcpServer:    NewServer(),
		Conn:         conn,
		ConnID:       0,
		MsgHandler:   client.GetMsgHandler(),
		ExitBuffChan: make(chan struct{}, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, sutils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
		onConnStart:  client.GetOnConnStart(),
		onConnStop:   client.GetOnConnStop(),
		wg:           &sync.WaitGroup{},
	}
	c.IsClosed.Store(false)
	c.IsAborted.Store(false)

	return c
}

func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is  running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		dp := NewDataPack()

		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConn(), headData); err != nil {
			c.IsClosed.Store(true)
			if err == io.EOF {
				fmt.Println("Connection closed by peer")
			} else {
				fmt.Println("read msg head error ", err)
				c.IsAborted.Store(true)
			}
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			c.IsClosed.Store(true)
			c.IsAborted.Store(true)
			fmt.Println("unpack error ", err)
			break
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConn(), data); err != nil {
				c.IsClosed.Store(true)
				c.IsAborted.Store(true)
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)

		req := Request{
			conn: c,
			msg:  msg,
			wg:   c.wg,
		}
		
		c.wg.Add(1)
		if sutils.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

func (c *Connection) StartWriter() {

	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Writer exit!]")

	for {
		select {
		case data, ok := <-c.msgChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Data error:, ", err, " Conn Writer exit")
				}
			} else {
				fmt.Println("msgChan is Closed")
				return
			}
		case data, ok := <-c.msgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("Send Buff Data error:, ", err, " Conn Writer exit")
				}
			} else {
				fmt.Println("msgBuffChan is Closed")
				return
			}
		}
	}
}

func (c *Connection) Start() {

	go c.StartWriter()
	
    go c.StartReader()

	c.onConnStart(c)

}

func (c *Connection) Stop() {
	if c.IsClosedOnce.Load() {
		return
	}

	c.IsClosedOnce.Store(true)
	c.IsClosed.Store(true)

	c.wg.Wait()

	c.onConnStop(c)
	c.Conn.Close()

	c.TcpServer.GetConnMgr().Remove(c)
	
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgId)
		return errors.New("Pack error msg ")
	}

	c.msgBuffChan <- msg

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, bool) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, true
	} else {
		return nil, false
	}
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}