package ziface

type IDataPack interface {
	GetHeadLen()                       // 获取包头长度的方法
	Pack(msg IMessage) ([]byte, error) // 封包的方法
	UnPack([]byte) (IMessage, error)   // 解包的方法
}
