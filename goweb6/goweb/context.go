package goweb

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/dengjiawen8955/go_utils/stringu"
)

const (
	ContentLength   = "content-length"
	ContentType     = "content-type"
	BufSize         = 1024
	JsonContentType = "application/json"
	WsMagicKeyPost  = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

type Context struct {
	Conn           net.Conn
	Addr           net.Addr          //conn 拿到的地址
	Method         string            //GET
	Path           string            //ping
	Proto          string            //HTTP/1.1
	Headers        map[string]string //请求头
	Body           []byte            //请求体数据啥的，直接放这里.
	Handlers       []HttpHandler     //中间件+业务handler
	HandlerIndex   int               //index 记录当前执行到第几个中间件
	StartTimeStamp int64             //会话开始的时间戳,用于计算耗时.
}

//GET /ping HTTP/1.1
//Host    img.mukewang.com
func newContext(conn net.Conn) (*Context, error) {
	ctx := &Context{
		Conn:           conn,
		Addr:           conn.RemoteAddr(),
		Headers:        make(map[string]string),
		HandlerIndex:   -1,
		StartTimeStamp: time.Now().UnixNano(),
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
	} else {
		return nil, fmt.Errorf("error: %s ", "http parse")
	}
	//1.2 其他 header 存起来.
	for {
		line, _, _ := reader.ReadLine()
		//log.Printf("%s\n", string(line))
		if len(line) == 0 {
			break
		}
		s := string(line)
		fmt.Printf("%s\n", s)
		split := strings.Split(s, ": ")
		if len(split) == 2 {
			ctx.Headers[strings.ToLower(split[0])] = split[1]
		} else {
			break
		}
	}
	//2. 请求体
	s = ctx.Headers[ContentLength]
	size, err := strconv.Atoi(s)
	if err != nil {
		size = 0
	}
	if size > 0 {
		data := make([]byte, size)
		conn.Read(data)
		ctx.Body = data
	}
	return ctx, nil
}
func (c *Context) Next() {
	c.HandlerIndex++
	s := len(c.Handlers)
	for ; c.HandlerIndex < s; c.HandlerIndex++ {
		c.Handlers[c.HandlerIndex](c)
	}
}

//Write json back and return.
func (c *Context) Json(data []byte) {
	c.writeHeader(JsonContentType, len(data))
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

//---------------- form data ------------

//FromData
type FromData struct {
	Name        string
	FileName    string
	ContentType string
	Data        []byte
}

//FormData 我
func (c *Context) GetFormData() (map[string]*FromData, error) {
	boundary := ""
	formDataMap := make(map[string]*FromData)
	//1. 拿到分割符
	ct, ok := c.Headers[ContentType]
	if !ok {
		return nil, fmt.Errorf("[error]:do not have contentType header")
	}
	splits := strings.Split(ct, "boundary=")
	if len(splits) != 2 {
		return nil, fmt.Errorf("[error]:boundary split err")
	}
	//赋值 boundary.
	boundary = splits[1]

	//2. 拿到数据
	buffer := bytes.NewBuffer(c.Body)
	for {
		//换行读取
		line, err := buffer.ReadBytes('\n')
		if err != nil {
			_, err := buffer.ReadBytes('\n')
			if err == io.EOF {
				break
			} else {
				continue
			}
		}
		//判断是否为分割符
		if string(line) == boundary {
			fromD := &FromData{}
			//读取接下
			line, err := buffer.ReadBytes('\n')
			if err != nil {
				break
			}
			str := string(line)
			reg := `Content-Disposition: (\S*?); name="(\S*?)"; filename="(\S*?)"`
			subs := stringu.GetSubStringByRegex(str, reg)
			fromD.Name = subs[1]
			fromD.FileName = subs[2]
			//读取接下
			line, err = buffer.ReadBytes('\n')
			if err != nil {
				break
			}
			str2 := string(line)
			regx2 := `Content-Type: (\S*)`
			subs2 := stringu.GetSubStringByRegex(str2, regx2)
			fromD.ContentType = subs2[0]
		}

	}

	return formDataMap, nil
}
