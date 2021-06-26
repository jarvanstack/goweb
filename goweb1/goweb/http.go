package goweb

import (
	"fmt"
	"log"
	"net"
)

type Http struct {
	Network string
	Addr    string
	r *router
}



func  newHttp( network ,addr string,r *router)*Http {
	return &Http{Network: network,Addr: addr,r: r}
}

func (h *Http)runHTTP()  {
	var err  error
	listen, err := net.Listen(h.Network, h.Addr)
	if err != nil {
		log.Printf("port already in use=%#v\n", h.Addr)
		log.Printf("err=%#v\n", err)
		return
	}
	for  {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("err=%#v\n", err)
			log.Printf("%s\n", "Accept() error")
			conn.Close()
			continue
		}
		log.Printf("conn.RemoteAddr()=%#v\n", conn.RemoteAddr())
		ctx, err := newContext(conn)
		if err != nil || ctx == nil{
			log.Printf("%s\n", "newContext(conn) err")
			conn.Close()
			continue
		}
		key := fmt.Sprintf(keyFmt,ctx.Method, ctx.Path)
		hf,ok := h.r.httpMap[key]
		if !ok {
			conn.Close()
			continue
		}
		//开一个协程去处理 http.
		//go func(ctx *Context) {
		//	hf(ctx)
		//}(ctx)
		hf(ctx)
	}
}
