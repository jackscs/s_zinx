package ziface

type IServer interface {
	Start()

	Stop()

	Serve()

	//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
	AddRouter(msgID uint32, router IRouter)

	GetConnMgr() IConnManger

	// 设置hook函数
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func(IConnection))

	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(IConnection))

	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)

	//该Server的连接断开时的Hook函数
	CallOnConnStop(conn IConnection)
}
