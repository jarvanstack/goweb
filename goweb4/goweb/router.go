package goweb

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restfulu"
	"log"
	"net"
	"strings"
)

const (
	keyFmt = "%s-%s" //比如 GET-/user/name
)

type HttpHandler func(ctx *Context)
type router struct {
	handlers    map[string]HttpHandler
	roots       map[string]*node
	contextPath string
	web *Web
}

func newRouter(contextPath string,web *Web) *router {
	return &router{
		handlers:    make(map[string]HttpHandler),
		roots:       make(map[string]*node),
		contextPath: contextPath,
		web: web,
	}
}
//添加到路由
//method: GET
//cp: = /bmft
//key: GET-/bmft/v1/ping
func (r *router) addRouter(method, urlPath string, httpHandler HttpHandler) {
	// GET-/bmft
	rootKey := fmt.Sprintf(keyFmt, method, r.contextPath)
	_, ok := r.roots[rootKey]
	if !ok {
		r.roots[rootKey] = newNode(rootKey)
	}
	//放入路由
	key := rootKey+urlPath
	r.roots[rootKey].insert(key,parsePath(urlPath),0)
	//放入handler map
	r.handlers[key] = httpHandler
}


func (r *router) runHTTP(network, addr string) {
	log.Printf("network=%s,addr%s\n", network, addr)
	var err error
	listen, err := net.Listen(network, addr)
	if err != nil {
		log.Printf("port already in use=%#v\n", addr)
		log.Printf("err=%#v\n", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("err=%#v\n", err)
			log.Printf("%s\n", "Accept() error")
			conn.Close()
			continue
		}
		go func(conn net.Conn) {
			ctx, err := newContext(conn)

			if err != nil || ctx == nil {
				log.Printf("%s\n", "newContext(conn) err")
				conn.Close()
				return
			}
			//注册 middleware 到 ctx
			var middlewares []HttpHandler
			urls := strings.Split(ctx.Path, r.contextPath)
			if len(urls) != 2 {
				ctx.Json(restfulu.NotFound("NOT_FOUND"))
				return
			}
			for _, group := range r.web.groups {
				if strings.HasPrefix(urls[1], group.prefix) {
					middlewares = append(middlewares, group.middlewares...)
				}
			}
			ctx.handlers = middlewares
			//拿到路由
			paths := parsePath(ctx.Path)
			rootKey := fmt.Sprintf(keyFmt,ctx.Method,"/"+paths[0])
			root := r.roots[rootKey]
			if root == nil {
				ctx.Json(restfulu.NotFound("NOT_FOUND"))
			}else {
				n := root.search(paths[1:])
				if n == nil {
					ctx.Json(restfulu.NotFound("NOT_FOUND"))
				}else {
					handler := r.handlers[n.keyOfHandler]
					//执行业务方法
					//业务方法也可以看做中间件，最后一个中间件
					ctx.handlers = append(ctx.handlers, handler)
					ctx.Next()
				}
			}
		}(conn)
	}
}
