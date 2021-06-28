package goweb

import (
	"bufio"
	"bytes"
	"net"
	"strconv"
	"strings"
)

const (
	contentLength   = "content-length"
	jsonContentTyp = "application/json"
)

type Context struct {
	Conn    net.Conn
	Addr    net.Addr          //conn 拿到的地址
	Method  string            //GET
	Path    string            //ping
	Proto   string            //HTTP/1.1
	Headers map[string]string //请求头
	body    []byte            //请求体数据啥的，直接放这里.


}



//GET /ping HTTP/1.1
//Host    img.mukewang.com
func newContext(conn net.Conn) (*Context, error) {
	ctx := &Context{Conn: conn, Addr: conn.RemoteAddr(), Headers: make(map[string]string)}
	//1.get headers
	reader := bufio.NewReader(conn)
	//1.1 请求方法，请求路径，请求协议.
	line, _, err := reader.ReadLine()
	s := string(line)
	if err != nil {
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

//Write json back and return.
func (c *Context) Json(data []byte) {
	defer c.Conn.Close()
	b := bytes.Buffer{}
	b.WriteString("HTTP/1.1 200 OK")
	b.WriteString("\r\n")
	b.WriteString("content-length: ")
	b.WriteString(strconv.Itoa(len(data)))
	b.WriteString("\r\n")
	b.WriteString("content-type: ")
	b.WriteString(jsonContentTyp)
	b.WriteString("; charset=utf-8")
	b.WriteString("\r\n")
	b.WriteString("\r\n")
	b.Write(data)
	c.Conn.Write(b.Bytes())
}
