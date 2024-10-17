package znet

type Message struct {
	Id      uint32
	DataLen uint32
	Data    []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (i *Message) GetDataLen() uint32 {
	return i.DataLen
}

func (i *Message) GetMsgID() uint32 {
	return i.Id
}

func (i *Message) GetData() []byte {
	return i.Data
}

func (i *Message) SetDataLen(len uint32) {
	i.DataLen = len
}

func (i *Message) SetMsgId(id uint32) {
	i.Id = id
}

func (i *Message) SetData(data []byte) {
	i.Data = data
}
