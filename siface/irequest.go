package siface

type IRequest interface {
	GetConn()       IConn
	GetData()       []byte
	GetMsgID()      uint32
	Done()
}
