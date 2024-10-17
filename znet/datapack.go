package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) HeadDataLen() uint32 {
	//Id uint32(4字节) +  DataLen uint32(4字节)
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})

	// 写入dataLen
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写入msgID
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgID()); err != nil {
		return nil, err
	}

	// 写入data
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil

}

func (dp *DataPack) UnPack(binnaryData []byte) (ziface.IMessage, error) {
	dataBuf := bytes.NewBuffer(binnaryData)

	msg := &Message{}

	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读取msgID
	if err := binary.Read(dataBuf, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen长度是否超过我们允许的最大长度
	if utils.GLOBALOBJ.MaxPacketSize > 0 && msg.DataLen > utils.GLOBALOBJ.MaxPacketSize {
		return nil, errors.New("to long message recv")
	}

	return msg, nil

}
