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
	// 请求，直接回复响应
	fmt.Println("[Sinx] Received:", string(request.GetData()))
	err := request.GetConn().SendMsg(0, []byte("received"))
	if err != nil {
		fmt.Println("[Sinx] error:", err)
	}
}

func main() {
	s := snet.NewServer()

	s.AddRouter(0, &helloRouter{})

	s.Start()

	time.Sleep(5 * time.Second)
	s.Stop()
}