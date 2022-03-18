package znet

import (
	"ZinxV0.2/src/zinx/utils"
	"ZinxV0.2/src/zinx/ziface"
	"bytes"
	"encoding/binary"
	"errors"
)

/*
	封包 拆包 模块
*/

type DataPack struct {
}

//拆包封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包的头长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	//Datalen uint32(4字节) +  ID uint32(4字节)
	return 8
}

//封包方法
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建存放byte字节流的缓冲
	databuffer := bytes.NewBuffer([]byte{})
	//将datalen 写入 databuffer
	if err := binary.Write(databuffer, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}
	//将MsgId 写入 databuffer
	if err := binary.Write(databuffer, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//将Data数据 写入 databuffer
	if err := binary.Write(databuffer, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return databuffer.Bytes(), nil
}

//拆包方法 (将包的Head信息读出) 再根据head信息中的data长度再读一次
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	databuffer := bytes.NewBuffer(binaryData)
	//只解压head信息,得到datalen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(databuffer, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读MsgID
	if err := binary.Read(databuffer, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断datalen是否已经超出了我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil
}
