package znet

import "zinx/ziface"

type Request struct {
	Conn    ziface.IConnection
	Message ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Message.GetData()
}

// 获取请求的消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.Message.GetMsgID()
}
