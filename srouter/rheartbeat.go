package srouter

import (
	"fmt"
	"time"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/snet"
)

const (
	MsgIDHeartbeatRequest  = 1001 // 心跳请求
	MsgIDHeartbeatResponse = 1002 // 心跳响应
)

type HeartbeatPingRouter struct {
	snet.BaseRouter
}

func (h *HeartbeatPingRouter) Handle(request siface.IRequest) {
	// 收到心跳请求，直接回复响应
	fmt.Println("[Sinx] Received heartbeat request, sending pong...")
	err := request.GetConn().SendMsg(MsgIDHeartbeatResponse, []byte("pong"))
	if err != nil {
		fmt.Println("[Sinx] Send heartbeat response error:", err)
	}
}

type HeartbeatPongRouter struct {
	snet.BaseRouter
}

func (h *HeartbeatPongRouter) Handle(request siface.IRequest) {
	// 收到心跳请求，直接回复响应
	request.GetConn().SetProperty("lastActiveTime", time.Now())
	fmt.Println("[Sinx] Received heartbeat response, updating last active time...")
}