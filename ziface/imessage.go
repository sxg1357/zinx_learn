package ziface

type IMessage interface {
	SetData([]byte)
	SetDataLen(uint32)
	SetMsgId(uint32)

	GetData() []byte
	GetDataLen() uint32
	GetMsgId() uint32
}
