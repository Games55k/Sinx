package siface

/*把客户端请求的链接信息和请求的数据包装到Request里*/
type IRequest interface {
	GetConn()       IConn
	GetData()       []byte
	GetMsgID()      uint32
}
