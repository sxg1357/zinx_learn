package znet

import "zinx_learn/ziface"

type Message struct {
	Data   []byte
	MsgLen uint32
	MsgId  uint32
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetDataLen(len uint32) {
	m.MsgLen = len
}

func (m *Message) GetDataLen() uint32 {
	return m.MsgLen
}

func (m *Message) SetMsgId(id uint32) {
	m.MsgId = id
}

func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}

func NewMessage(data []byte, len uint32, msgId uint32) ziface.IMessage {
	return &Message{
		Data:   data,
		MsgLen: len,
		MsgId:  msgId,
	}
}
