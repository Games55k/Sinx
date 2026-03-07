package sutils

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/Games55k/Sinx/siface"
)

type GlobalObj struct {
	TcpServer     siface.IServer // 当前Zinx的全局Server对象
	Host          string         // 当前服务器主机IP
	TcpPort       int            // 当前服务器主机监听端口号
	Name          string         // 当前服务器名称
	Version       string         // 当前Zinx版本号
	MaxPacketSize uint32         // 允许的最大包长度
	MaxConn       int            // 当前服务器主机允许的最大连接个数
}

var GlobalObject *GlobalObj

// 读取用户配置文件。
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/Sinx.json")
	if err != nil {
		// 在没有配置文件时保留默认配置，避免测试/库加载时直接panic。
		if errors.Is(err, os.ErrNotExist) {
			return
		}
		panic(err)
	}

	if err := json.Unmarshal(data, g); err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "SinxServerApp",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}

	GlobalObject.Reload()
}
