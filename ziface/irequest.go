package ziface

type IRequest interface {
	GetData() []byte
	GetConnection() IConnection
}
