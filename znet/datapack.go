package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zinx_learn/utils"
	"zinx_learn/ziface"
)

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	bufferData := bytes.NewBuffer([]byte{})

	//打包数据长度
	if err := binary.Write(bufferData, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//打包消息id
	if err := binary.Write(bufferData, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//打包数据内容
	if err := binary.Write(bufferData, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return bufferData.Bytes(), nil
}

func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}

	//解析数据长度
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgLen); err != nil {
		return nil, err
	}

	if msg.MsgLen > utils.GlobalObject.MaxPacketLen {
		return nil, errors.New("exceed max package len")
	}

	//解析消息id
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.MsgId); err != nil {
		return nil, err
	}
	return msg, nil
}
