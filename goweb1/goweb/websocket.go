package goweb

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/base_util"
	"net"
)



type WsContext struct {
	Conn net.Conn
}

func newWs(conn net.Conn) *WsContext {
	return &WsContext{Conn: conn}
}
//1. 首读取 2 字节，
//2. 拿到数据长度
//3. 判断数据长度类型并拿到数据长度类型，和掩码
//4. 使用掩码和数据长度数据拿到最后的数据
func (w *WsContext) ReadMsg() ([]byte,error){
	//1.
	var err error
	conn := w.Conn
	h := make([]byte,2)
	conn.Read(h)
	//2.
	// & 将第一位变为 0
	dataLen := 0
	mask := make([]byte,4)
	l := h[1] & 127
	//3.
	if l == 126 {
		next := make([]byte,6)
		conn.Read(next)
		dataLen,err = base_util.BytesToInt(next[:2],false)
		if	err!=nil {
			return nil, err
		}
		copy(mask,next[2:6])
	}else if l == 127 {
		next := make([]byte,10)
		conn.Read(next)
		dataLen,err = base_util.BytesToInt(next[:8],false)
		if	err!=nil {
			return nil, err
		}
		copy(mask,next[8:10])
	}else {
		dataLen = int(l)
		next := make([]byte, 4)
		conn.Read(next)
		copy(mask,next)
	}
	//4.
	data := make([]byte,dataLen)
	conn.Read(data)
	for i, b := range data {
		data[i] = b ^ mask[(i%4)]
	}
	return data,nil
}
//写入数据
//1. 拿到 bytes 长度
//2. 通过长度判断并获得总长度并写入协议头
//3. 数据拷贝载荷
func (w *WsContext) WriteMsg(data []byte) error {
	var err error
	dataL := len(data)
	payloadSize := 0
	conn := w.Conn
	if dataL < 126{
	}else if dataL >= 126 && dataL < int(base_util.Uint8Max) {
		payloadSize = 2
	}else  {
		payloadSize = 8
	}
	//帧
	frame := make([]byte,2+payloadSize+dataL)
	//写入第一byte 1000 0001 128 + 1 = 129
	b1, err := base_util.UintToBytes(129, 1)
	if err != nil {
		return err
	}
	frame[0] = b1[0]

	//写入额外的play len
	switch payloadSize {
	case 0:
		//写入第二帧长度
		b2, err := base_util.IntToBytes(dataL,1)
		if err != nil {
			return err
		}
		frame[1] = b2[0]
		break
	case 2:
		//写入第二帧长度
		b2, err := base_util.IntToBytes(126,1)
		if err != nil {
			return err
		}
		frame[1] = b2[0]

		break
	case 8:
		//写入第二帧长度
		b2, err := base_util.IntToBytes(127,1)
		if err != nil {
			return err
		}
		frame[1] = b2[0]
		break
	default:
		return fmt.Errorf("error:%s", "Not have this payloadSize")
	}
	if payloadSize>0 {
		//写入额外长度
		b3, err := base_util.UintToBytes(uint64(dataL), uint8(payloadSize))
		if err != nil {
			return err
		}
		copy(frame[2:2+payloadSize],b3[:payloadSize])
	}
	//写入数据
	copy(frame[2+payloadSize:2+payloadSize+dataL],data[:dataL])
	//回写客户端值
	_, err = conn.Write(frame)
	return err
}