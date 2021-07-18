package goweb

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/dengjiawen8955/go_utils/stringu"
	"net"
	"strconv"
	"strings"
)

const (
	contentLength   = "content-length"
	bufSize         = 1024
	jsonContentType = "application/json"
	WsMagicKeyPost  = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

type Context struct {
	Conn     net.Conn
	Addr     net.Addr          //conn 拿到的地址
	Method   string            //GET
	Path     string            //ping
	Proto    string            //HTTP/1.1
	Headers  map[string]string //请求头
	body     []byte            //请求体数据啥的，直接放这里.
	handlers []HttpHandler     //中间件+业务handler
	index    int               //index 记录当前执行到第几个中间件
}

//GET /ping HTTP/1.1
//Host    img.mukewang.com
func newContext(conn net.Conn) (*Context, error) {
	ctx := &Context{
		Conn: conn,
		Addr: conn.RemoteAddr(),
		Headers: make(map[string]string),
		index:  -1,
	}
	//1.GET headers
	reader := bufio.NewReader(conn)
	//1.1 请求方法，请求路径，请求协议.
	line, _, err := reader.ReadLine()
	//log.Printf("%s\n", string(line))
	s := string(line)
	if err != nil || strings.EqualFold(s, "\r\n") {
		//HTTP 头解析错误
		return nil, err
	}
	split := strings.Split(s, " ")
	if len(split) == 3 {
		ctx.Method = split[0]
		ctx.Path = split[1]
		ctx.Proto = split[2]
	}
	//1.2 其他 header 存起来.
	for {
		line, _, _ := reader.ReadLine()
		//log.Printf("%s\n", string(line))
		if len(line) == 0 {
			break
		}
		s := string(line)
		split := strings.Split(s, ": ")
		if len(split) == 2 {
			ctx.Headers[split[0]] = split[1]
		} else {
			break
		}
	}
	//2. 请求体
	s = ctx.Headers[contentLength]
	size, err := strconv.Atoi(s)
	if err != nil {
		size = 0
	}
	if size > 0 {
		data := make([]byte, size)
		rn, _ := conn.Read(data)
		ctx.body = data[:rn]
	}
	return ctx, nil
}
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

//Write json back and return.
func (c *Context) Json(data []byte) {
	c.writeHeader(jsonContentType, len(data))
	c.writeBody(data)
	c.Conn.Close()
}
func (c *Context) writeHeader(contentType string, bodySize int) (int, error) {
	conn := c.Conn
	buffer := bytes.Buffer{}
	buffer.WriteString("HTTP/1.1 200 OK\r\n")
	buffer.WriteString("content-length: ")
	buffer.WriteString(strconv.Itoa(bodySize))
	buffer.WriteString("\r\n")
	buffer.WriteString("content-type: ")
	buffer.WriteString(contentType)
	buffer.WriteString(";charset=utf-8\r\n\r\n")
	return conn.Write(buffer.Bytes())
}
func (c *Context) writeBody(data []byte) (int, error) {
	conn := c.Conn
	return conn.Write(data)
}

//1.创建 ws 对象
// 1.1. 发送协议头
// 1.2. 创建 ws 对象.
//2. 使用 ws 对象持续通信.
func (c *Context) NewWs() (*WsContext, error) {
	var err error
	//1.1
	key, ok := c.Headers["Sec-WebSocket-Key"]
	if !ok {
		err = fmt.Errorf("error:%s", "not have Sec-WebSocket-Key error")
		return nil, err
	}
	conn := c.Conn
	key = key + WsMagicKeyPost
	s1 := stringu.GetSha1ByStr(key)
	md := stringu.GetMd5ByBytes(s1)
	b := bytes.Buffer{}
	b.WriteString("HTTP/1.1 101 Switching Protocols\r\n")
	b.WriteString("Upgrade: websocket\r\n")
	b.WriteString("Connection: Upgrade\r\n")
	b.WriteString("Sec-WebSocket-Accept: ")
	b.WriteString(md)
	b.WriteString("\r\n\r\n")
	conn.Write(b.Bytes())
	//1.2
	return newWs(conn), nil
}
