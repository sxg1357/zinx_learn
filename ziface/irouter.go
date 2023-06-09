package ziface

type IRouter interface {
	PreHandler(request IRequest)
	Handler(request IRequest)
	PostHandler(request IRequest)
}
