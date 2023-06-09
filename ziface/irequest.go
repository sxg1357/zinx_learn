package ziface

type IRequest interface {
	GetData() []byte
	GetConnection() IConnection
	GetMsgId() uint32
}
