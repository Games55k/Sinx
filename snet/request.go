package snet

import (
	"sync"
	"github.com/Games55k/Sinx/siface"
)

type Request struct {
	conn siface.IConn
	msg  siface.IMessage
	wg   *sync.WaitGroup
}

func (r *Request) GetConn() siface.IConn {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) Done() {
	r.wg.Done()
}