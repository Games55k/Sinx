package main

import (
	"fmt"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/snet"
)

type PingRouter struct {
	snet.BaseRouter
}

func (this *PingRouter) Handle(request siface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConn().SendBuffMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloSinxRouter struct {
	snet.BaseRouter
}

func (this *HelloSinxRouter) Handle(request siface.IRequest) {
	fmt.Println("Call HelloSinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConn().SendBuffMsg(1, []byte("Hello Sinx Router"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn siface.IConn) {
	fmt.Println("DoConnecionBegin is Called ... ")
	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn siface.IConn) {
	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {
	s := snet.NewServer()

	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloSinxRouter{})

	s.Serve()
}