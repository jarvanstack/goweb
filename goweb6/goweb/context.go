package goweb

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dengjiawen8955/go_utils/stringu"
)

//常量
const (
	ContentLength   = "content-length"
	ContentType     = "content-type"
	BufSize         = 1024
	JsonContentType = "application/json"
	WsMagicKeyPost  = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

//变量
var once sync.Once

type Context struct {
	Conn    net.Conn
	Addr    net.Addr          //conn 拿到的地址
	Method  string            //GET
	Path    string            //ping
	Proto   string            //HTTP/1.1
	Headers map[string]string //请求头
	body    []byte            //异步请求体 注意:使用的时候请使用 GetBody() 方式来获取

	Handlers       []HttpHandler //中间件+业务handler
	HandlerIndex   int           //index 记录当前执行到第几个中间件
	StartTimeStamp int64         //会话开始的时间戳,用于计算耗时.
	BodyReady      chan int      //管道读取 body body 数据读取成功后会进入会再管道中存放值
	BodySize       int           //数据体的大小
}

//解析前端 Json 数据, 需要传入 struct 指针.
//Example:
// 	u := &U{}
// 	ctx.Unmarshal(u)
// 	ctx.Json(restfulu.Ok(u.Name))
func (c *Context) Unmarshal(v interface{}) error {
	if !strings.EqualFold(c.Headers[ContentType], JsonContentType) {
		return fmt.Errorf("[error]:content-type not %s", JsonContentType)
	}
	err := json.Unmarshal(c.GetBody(), v)
	if err != nil {
		return fmt.Errorf("[error]:json Unmarshal err")
	}
	return nil
}

// 开发阶段,不能使用
//TODO: 开发阶段
func (c *Context) GetForm() (*Form, error) {
	boundary := ""
	endBoundary := ""
	isFinish := false
	form := &Form{
		FormFileMap: make(map[string]*FormFile),
		FormDataMap: make(map[string]*FormData),
	}
	// 1. 拿到分割符
	ct, ok := c.Headers[ContentType]
	if !ok {
		return nil, fmt.Errorf("[error]:do not have contentType header")
	}
	splits := strings.Split(ct, "boundary=")
	if len(splits) != 2 {
		return nil, fmt.Errorf("[error]:boundary split err")
	}
	//赋值 boundary.
	boundary = "--" + splits[1] + "\r\n"
	endBoundary = "--" + splits[1] + "--\r\n"
	//拿到数据body
	buffer := bytes.NewBuffer(c.GetBody())
	//取出第一行没用的分隔符
	_, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("[error]:第一个boundary解析失败")
	}
	for {
		//2.file 或者 data
		l2, err := buffer.ReadBytes('\n')
		if err != nil {
			break
		}
		splits := strings.Split(string(l2), "; ")
		if len(splits) == 3 {
			file := &FormFile{}
			//file
			reg := `name="(\S*?)"`
			subs, err := stringu.GetSubStringByRegex(splits[1], reg)
			if err != nil {
				continue
			}
			file.Name = subs[0]
			reg = `filename="(\S*?)"`
			subs, err = stringu.GetSubStringByRegex(splits[2], reg)
			if err != nil {
				continue
			}
			file.FileName = subs[0]
			line, err := buffer.ReadBytes('\n')
			if err != nil {
				continue
			}
			reg = `Content-Type: (\S*?)`
			subs, err = stringu.GetSubStringByRegex(string(line), reg)
			if err != nil {
				continue
			}
			file.ContentType = subs[0]
			//开始读取数据
			var data bytes.Buffer
			for {
				line, err := buffer.ReadBytes('\n')
				s := string(line)
				if err != nil {
					break
				}
				if strings.Compare(s, boundary) == 0 {
					break
				}
				//结束
				if strings.Compare(s, endBoundary) == 0 {
					isFinish = true
					break
				}
				//内存拷贝
				_, err = data.Write(line)
				if err != nil {
					//内存爆炸才会抛出异常.
					return nil, fmt.Errorf("[error]:buffer becomes too large, Write will panic with ErrTooLarge")
				}
			}
			bs := data.Bytes()
			file.Data = bs[4 : len(bs)-2]
			form.FormFileMap[file.Name] = file
			if isFinish {
				return form, nil
			}
		} else if len(splits) == 2 {
			//data
			formD := &FormData{}
			//file
			reg := `name="(\S*?)"`
			subs, err := stringu.GetSubStringByRegex(splits[1], reg)
			if err != nil {
				continue
			}
			formD.Name = subs[0]
			//开始读取数据
			var data bytes.Buffer
			for {
				line, err := buffer.ReadBytes('\n')
				s := string(line)
				if err != nil {
					break
				}
				if strings.Compare(s, boundary) == 0 {
					break
				}
				//结束
				if strings.Compare(s, endBoundary) == 0 {
					isFinish = true
					break
				}
				//内存拷贝
				_, err = data.Write(line)
				if err != nil {
					break
				}
			}
			bs := data.Bytes()
			formD.Data = bs[4 : len(bs)-2]
			form.FormDataMap[formD.Name] = formD
			if isFinish {
				return form, nil
			}
		}
	}
	return nil, fmt.Errorf("[error]:FormData第二行数据解析失败")
}

func (c *Context) GetBody() []byte {
	// <-c.BodyReady从管道中取出值
	// 发送到管道 c.BodyReady <- 1
	once.Do(
		func() {
			if c.BodySize > 0 {
				// fmt.Printf("%s\n", "等数据发送过来")
				<-c.BodyReady
				//等数据发送过来
			}
		})
	return c.body
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
		BodyReady:      make(chan int, 1), // 添加一个缓冲,让那个线程快点退出.
	}
	//1.GET headers
	reader := bufio.NewReader(conn)
	//1.1 请求方法，请求路径，请求协议.
	line, _, err := reader.ReadLine()
	s := string(line)
	fmt.Printf("%s\n", s)
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
	ctx.BodySize = size
	if err != nil {
		size = 0
	}
	if size > 0 {
		go func() {
			//TODO: 为什么这里开不了协程??? 一开就读取不了数据?
			//不要使用 conn conn 的权限给 Reader 了,现在只有 reader
			data := make([]byte, size)
			reader.Read(data)
			ctx.body = data
			// fmt.Printf("%s\n", string(ctx.body))
			ctx.BodyReady <- 1
		}()

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

type Form struct {
	FormFileMap map[string]*FormFile //key 为 name
	FormDataMap map[string]*FormData //key 为 name
}

// Form 文件类型
// Example
// 	Content-Disposition: form-data; name="uploadFile"; filename="f1.txt"
// 	Content-Type: text/plain
//
//	文件1
type FormFile struct {
	Name        string
	FileName    string
	ContentType string
	Data        []byte
}

// Form 表单键值对类型
// Example
//	Content-Disposition: form-data; name="sign_time"
//
//	2019
type FormData struct {
	Name string
	Data []byte
}
