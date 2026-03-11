package snet

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Games55k/Sinx/siface"
)

type ConnManager struct {
	conns      sync.Map
	connsCount atomic.Int32
	isClosed   atomic.Bool
}

func NewConnManager() *ConnManager {
	return &ConnManager{}
}

func (connMgr *ConnManager) Add(conn siface.IConn) error {
	if connMgr.isClosed.Load() {
		return errors.New("ConnManager closed")
	}

	connMgr.conns.Store(conn.GetConnID(), conn)
	connMgr.connsCount.Add(1)

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
	return nil
}

func (connMgr *ConnManager) Remove(conn siface.IConn) {
	if _, loaded := connMgr.conns.LoadAndDelete(conn.GetConnID()); loaded {
		connMgr.connsCount.Add(-1)
		fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully: conn num = ", connMgr.Len())
	}
}

func (connMgr *ConnManager) Get(connID uint32) (siface.IConn, error) {
	if conn, ok := connMgr.conns.Load(connID); ok {
		return conn.(siface.IConn), nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return int(connMgr.connsCount.Load())
}

func (connMgr *ConnManager) ClearConns() {
	connMgr.isClosed.Store(true)

	connMgr.conns.Range(func(key, value any) bool {
		value.(siface.IConn).Stop()
		return true
	})

	fmt.Println("Clear All Conns successfully: conn num = ", connMgr.Len())
}
