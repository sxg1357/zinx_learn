package ziface

type IServer interface {
	Server()
	Start()
	Stop()
	AddRouter(route IRouter)
}
