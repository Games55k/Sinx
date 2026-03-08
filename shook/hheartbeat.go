package shook

import (
	"fmt"
	"time"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/snet"
	"github.com/Games55k/Sinx/srouter"
)

func StartHeartbeat(conn siface.IConn) {
	conn.SetProperty("lastActiveTime", time.Now())
	ticker := time.NewTicker(5 * time.Second) // 5秒发送一次心跳
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.SendBuffMsg(srouter.MsgIDHeartbeatRequest, []byte("ping"))
			if err != nil {
				fmt.Println("[Cinx] Send heartbeat request error:", err)
				return
			}
		case <-conn.(*snet.Connection).ExitChan:
			return
		}
	}
}

func StartHeartbeatChecker(conn siface.IConn) {
	// 启动一个协程定期检查连接活性
	for {
		time.Sleep(5 * time.Second)

		// 获取最后一次活跃时间
		lastActiveTime, ok := conn.GetProperty("lastActiveTime")
		if !ok {
			continue
		}

		// 超过 15 秒未收到消息，判定超时
		if time.Since(lastActiveTime.(time.Time)) > 15*time.Second {
			fmt.Println("Connection timeout, closing...")
			conn.SendBuffMsg(404, []byte("timeout"))
			time.Sleep(5 * time.Second)
			conn.Stop()
			return
		}
	}
}