package ziface

type IMessage interface {
	GetDataLen() uint32
	GetMsgID() uint32
	GetData() []byte

	SetMsgId(uint32)   //设计消息ID
	SetData([]byte)    //设计消息内容
	SetDataLen(uint32) //设置消息数据段长度
}
