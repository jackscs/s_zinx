package utils

import (
	"encoding/json"
	"os"
	"zinx/ziface"
)

const (
	name             = "ZinxServerApp"
	version          = "v0.9"
	tcpPort          = 7777
	host             = "0.0.0.0"
	maxConn          = 12000
	maxPackSize      = 4096
	configPath       = "conf/zinx.json"
	workerPoolSize   = 10
	maxWorkerTaskLen = 1024
	maxMsgChanLen    = 5
)

var GLOBALOBJ *GlobalObj

type GlobalObj struct {
	TcpServer     ziface.IServer //当前Zinx的全局Server对象
	Host          string
	TcpPort       int
	Name          string
	Version       string
	MaxPacketSize uint32
	MaxConn       int

	ConfigPath string

	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32
	MaxMsgChanLen    uint32
}

func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("./conf/zinx.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GLOBALOBJ)
	if err != nil {
		panic(err)
	}

}

func init() {

	GLOBALOBJ = &GlobalObj{
		Name:          name,
		Version:       version,
		TcpPort:       tcpPort,
		Host:          host,
		MaxConn:       maxConn,
		MaxPacketSize: maxPackSize,

		ConfigPath:       configPath,
		WorkerPoolSize:   workerPoolSize,
		MaxWorkerTaskLen: maxWorkerTaskLen,
		MaxMsgChanLen:    maxMsgChanLen,
	}

	GLOBALOBJ.Reload()
}
