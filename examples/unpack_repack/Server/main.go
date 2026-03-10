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

	err := request.GetConn().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := snet.NewServer()

	s.AddRouter(0, &PingRouter{})

	s.Serve()
}
