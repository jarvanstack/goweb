package goweb

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restfulu"
	"log"
	"net"
)

const (
	keyFmt = "%s-%s" //比如 GET-/user/name
)

type HttpHandler func(ctx *Context)
type router struct {
	handlers    map[string]HttpHandler
	roots       map[string]*node
	contextPath string
}

func newRouter(contextPath string) *router {
	return &router{
		handlers:    make(map[string]HttpHandler),
		roots:       make(map[string]*node),
		contextPath: contextPath,
	}
}
//添加到路由
//method: GET
//cp: = /bmft
//key: GET-/bmft/goweb/ping
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
					handler(ctx)
				}
			}
		}(conn)
	}
}
