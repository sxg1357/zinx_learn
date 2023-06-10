package ziface

type IServer interface {
	Server()
	Start()
	Stop()
	AddRouter(msgId uint32, route IRouter)
	GetConnMgr() IConnManager
}
