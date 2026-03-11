package snet

import (
	"fmt"
	"strconv"

	"github.com/Games55k/Sinx/siface"
	"github.com/Games55k/Sinx/sutils"
)

type MsgHandle struct {
	Apis           map[uint32]siface.IRouter
	WorkerPoolSize uint32
	TaskQueues     []chan siface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]siface.IRouter),
		WorkerPoolSize: sutils.GlobalObject.WorkerPoolSize,
		TaskQueues:     make([]chan siface.IRequest, sutils.GlobalObject.WorkerPoolSize),
	}
}

// 非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request siface.IRequest) {
	defer request.Done()

	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router siface.IRouter) {
	// 判断当前 msg 绑定的 API 处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("[Sinx] Repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	// 添加 msg 与 API 的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("[Sinx] Add api msgId = ", msgId)
}

// Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan siface.IRequest) {
	// 监听队列中的消息
	for request := range taskQueue {
		mh.DoMsgHandler(request)
	}
}

// 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 初始化当前 worker 消息队列管道
		mh.TaskQueues[i] = make(chan siface.IRequest, sutils.GlobalObject.MaxWorkerTaskLen)

		fmt.Println("[Sinx] Worker ID = ", i, " is started.")
		// 创建一个 worker 协程
		go mh.StartOneWorker(i, mh.TaskQueues[i])
	}
}

// 分发消息给消息队列处理
func (mh *MsgHandle) SendMsgToTaskQueue(request siface.IRequest) {
	// 朴素的任务分配策略
	workerID := request.GetConn().GetConnID() % mh.WorkerPoolSize
	fmt.Println("[Sinx] Add ConnID=", request.GetConn().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)

	// 将消息发送给对应的 worker 的消息队列
	mh.TaskQueues[workerID] <- request
}

func (mh *MsgHandle) Stop() {
	for _, c := range mh.TaskQueues {
		if c != nil {
			close(c)
		}
	}
}
