package main

import (
	"fmt"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/snet"
)

// ping test 自定义路由
type PingRouter struct {
	snet.BaseRouter //一定要先基础BaseRouter
}

func (this *PingRouter) PreHandle(request siface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (this *PingRouter) Handle(request siface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (this *PingRouter) PostHandle(request siface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {
	//创建一个server句柄
	s := snet.NewServer()

	s.AddRouter(&PingRouter{})

	//2 开启服务
	s.Serve()
}
