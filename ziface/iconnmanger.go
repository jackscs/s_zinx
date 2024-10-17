package ziface

type IConnManger interface {
	Add(conn IConnection)                       // 添加连接
	Remove(conn IConnection)                    // 删除连接
	GetConn(connID uint32) (IConnection, error) // 利用connID获取连接
	Len() int                                   // 获取当前连接个数
	ClearConn()                                 // 停止并删除所有连接
}
