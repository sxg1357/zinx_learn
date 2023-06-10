package znet

import (
	"fmt"
	"strconv"
	"zinx_learn/utils"
	"zinx_learn/ziface"
)

type IMsgHandler struct {
	Apis           map[uint32]ziface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan ziface.IRequest
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

func (ih *IMsgHandler) StartOneWorker(workerId int, taskQueue chan ziface.IRequest) {
	fmt.Println("WorkerId = ", workerId, "is started")
	for {
		select {
		case request := <-taskQueue:
			ih.DoMsgHandler(request)
		}
	}
}

func (ih *IMsgHandler) StartWorkerPool() {
	for i := 0; i < int(ih.WorkerPoolSize); i++ {
		ih.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go ih.StartOneWorker(i, ih.TaskQueue[i])
	}
}

func (ih *IMsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	WorkerId := request.GetConnection().GetConnId() % ih.WorkerPoolSize
	ih.TaskQueue[WorkerId] <- request
}

func NewMsgHandler() *IMsgHandler {
	return &IMsgHandler{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
