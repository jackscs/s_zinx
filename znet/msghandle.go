package znet

import (
	"fmt"
	"strconv"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	Apis           map[uint32]ziface.IRouter //存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                    //业务工作Worker池的数量
	TaskQueue      []chan ziface.IRequest    //Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GLOBALOBJ.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GLOBALOBJ.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsghandle(request ziface.IRequest) {
	handle, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 判断当前msg绑定的api方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgID)))
	}

	// 添加msg与api的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api msgId = ", msgID)
}

func (mh *MsgHandle) StartWorkerPool() {

	for i := 0; i < int(utils.GLOBALOBJ.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GLOBALOBJ.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}

}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workerID, " is started.")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsghandle(request)
		}
	}
}

// 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request ziface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	// 获取workerID
	workerID := request.GetConnection().GetConnID() % utils.GLOBALOBJ.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgID(), "to workerID=", workerID)

	// 将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request

}
