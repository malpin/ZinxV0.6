package znet

import (
	"fmt"
)

type Message struct {
	Id      uint32 //消息id
	DataLen uint32 //消息长度
	Data    []byte //消息内容
}

//创建Message消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

//获取消息的ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

//获取消息长度
func (m *Message) GetMsgLen() uint32 {
	return m.DataLen
}

//获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

//设置消息的ID
func (m *Message) SetMsgId(id uint32) {
	m.Id = id
}

//设置消息的长度
func (m *Message) SetMsgLen(len uint32) {
	m.DataLen = len
}

//设置消息的内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) ToString() {
	fmt.Println("消息的ID", m.Id)
	fmt.Println("消息的长度", m.DataLen)
	fmt.Println("消息的内容", string(m.Data))
}
