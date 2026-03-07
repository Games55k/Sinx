package main

import (
	"fmt"
	"time"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/snet"
)

type helloRouter struct {
	snet.BaseRouter
}

func (h *helloRouter) Handle(request siface.IRequest) {
	fmt.Println("[Cinx] Received:", string(request.GetData()))
}

func main() {

	fmt.Println("Client Test ... start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(1 * time.Second)

	client := snet.NewClient("Client", "tcp4", "127.0.0.1", 7777)
	client.AddRouter(0, &helloRouter{})
	client.Start()

	for {
		time.Sleep(3 * time.Second)
		client.Conn().SendMsg(0, []byte("hello"))
	}
}