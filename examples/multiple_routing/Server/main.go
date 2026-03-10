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

	err := request.GetConn().SendMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloCinxRouter struct {
	snet.BaseRouter
}

func (this *HelloCinxRouter) Handle(request siface.IRequest) {
	fmt.Println("Call HelloSinxRouter Handle")
	fmt.Println("recv from client : msgId=", request.GetMsgID(), ", data=", string(request.GetData()))

	err := request.GetConn().SendMsg(1, []byte("Hello Sinx Router"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := snet.NewServer()

	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloCinxRouter{})

	s.Serve()
}
