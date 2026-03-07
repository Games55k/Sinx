package snet

import "github.com/Games55k/Sinx/siface"

type Request struct {
	conn siface.IConnection //已经和客户端建立好的 链接
	data []byte             //客户端请求的数据
}

// 获取请求连接信息
func (r *Request) GetConnection() siface.IConnection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.data
}