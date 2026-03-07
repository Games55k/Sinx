package main

import (
	"fmt"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/snet"
	"github.com/Games55k/Sinx/srouter"
	"github.com/Games55k/Sinx/shook"
)

type helloRouter struct {
	snet.BaseRouter
}

func (h *helloRouter) Handle(request siface.IRequest) {
	// 请求，直接回复响应
	fmt.Println("[Sinx] Received:", string(request.GetData()))
	err := request.GetConnection().SendMsg(0, []byte("received"))
	if err != nil {
		fmt.Println("[Sinx] error:", err)
	}
}

func main() {
	//创建一个server句柄
	s := snet.NewServer()
	s.AddRouter(srouter.MsgIDHeartbeatRequest, &srouter.HeartbeatPingRouter{})
	s.AddRouter(srouter.MsgIDHeartbeatResponse, &srouter.HeartbeatPongRouter{})
	s.SetOnConnStart(func(conn siface.IConnection) {
		go shook.StartHeartbeat(conn)
		go shook.StartHeartbeatChecker(conn)
	})
	s.AddRouter(0, &helloRouter{})
	//开启服务
	s.Serve()
}