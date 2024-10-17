package ziface

import "net"

type IConnection interface {
	Start()

	Stop()

	GetConnID() uint32

	// 从当前连接中获得原始socket连接
	GetTcpConnection() *net.TCPConn

	// 获取远程客户端的地址
	RemoteAddr() net.Addr

	// 将数据发送给远程的客户端
	SendMsg(msgID uint32, data []byte) error

	// 将数据通过缓冲通道发送给客户端
	SendBuffMsg(msgID uint32, data []byte) error

	// 设置连接属性
	SetProperty(key string, value interface{}) error
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 去除连接属性
	RemoveProperty(key string)
}

type HandFunc func(*net.TCPConn, []byte, int) error
