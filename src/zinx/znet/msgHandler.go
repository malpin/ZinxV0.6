package znet

import (
	"ZinxV0.2/src/zinx/ziface"
	"fmt"
	"strconv"
)

/*
消息管理模块实现
*/

type MagHandle struct {
	//存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

//初始化/创建MagHandle方法
func NewMagHandle() *MagHandle {
	return &MagHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

//调度/执行对应的Router消息处理方法
func (mh *MagHandle) DoMsgHandler(request ziface.IRequest) {
	//1.从IRequest 得到msgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("id为,", request.GetMsgID(), "没有找到对应方法")
	}
	//根据msgID调用 对应的业务
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

//为消息添加具体的处理逻辑
func (mh *MagHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("此id有对应处理方法" + strconv.Itoa(int(msgID)))
	}
	//添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("添加路由处理方法")
}
