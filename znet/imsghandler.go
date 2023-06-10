package znet

import (
	"fmt"
	"strconv"
	"zinx_learn/ziface"
)

type IMsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func (ih *IMsgHandler) AddRouter(msgId uint32, router ziface.IRouter) {
	if _, ok := ih.Apis[msgId]; ok {
		panic("repeated api, msgid " + strconv.Itoa(int(msgId)))
	}
	ih.Apis[msgId] = router
	fmt.Println("add api msgid = ", msgId)
}

func (ih *IMsgHandler) DoMsgHandler(request ziface.IRequest) {
	handler, ok := ih.Apis[request.GetMsgId()]
	if !ok {
		panic("apis doesn't exist msgid" + strconv.Itoa(int(request.GetMsgId())))
	}
	handler.PreHandler(request)
	handler.Handler(request)
	handler.PostHandler(request)
}

func NewMsgHandler() *IMsgHandler {
	return &IMsgHandler{
		Apis: make(map[uint32]ziface.IRouter),
	}
}
