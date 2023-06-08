package ziface

type IServer interface {
	Server()
	Start()
	Stop()
}
