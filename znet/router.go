package znet

import "zinx_learn/ziface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandler(request ziface.IRequest)  {}
func (br *BaseRouter) Handler(request ziface.IRequest)     {}
func (br *BaseRouter) PostHandler(request ziface.IRequest) {}
