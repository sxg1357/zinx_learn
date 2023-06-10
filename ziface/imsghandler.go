package ziface

type IMsgHandler interface {
	AddRouter(msgId uint32, router IRouter)
	DoMsgHandler(request IRequest)
}
