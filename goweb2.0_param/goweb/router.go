package goweb

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restfulu"
	"log"
	"net"
	"strings"
)

const (
	keyFmt = "%s-%s" //比如 get-/user/name
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
//method
//cp = /bmft
//path = /v1/ping
func (r *router) addRouter(method, urlPath string, httpHandler HttpHandler) {
	// get-/bmft
	rootKey := fmt.Sprintf(keyFmt, method, r.contextPath)
	_, ok := r.roots[rootKey]
	if !ok {
		r.roots[rootKey] = &node{}
	}
	key := rootKey+urlPath
	r.roots[rootKey].insert(key,parsePattern(urlPath),0)
	//get-/bmft/v1/ping
	r.handlers[key] = httpHandler
}

//获取路由
//get, /bmft/:lang/ping
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	// bmft :lang ping
	searchParts := parsePattern(path)
	rootKey := fmt.Sprintf(keyFmt, method, r.contextPath)
	params := make(map[string]string)
	root, ok := r.roots[rootKey]
	if !ok {
		return nil, nil
	}
	n := root.search(searchParts[1:], 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
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
		//go func(conn net.Conn) {
			ctx, err := newContext(conn)
			if err != nil || ctx == nil {
				log.Printf("%s\n", "newContext(conn) err")
				conn.Close()
				return
			}
			n, params := r.getRoute(ctx.Method, ctx.Path)
			fmt.Printf("n=%#v\n", n)
			fmt.Printf("params=%#v\n", params)
			if n != nil {
				ctx.Params = params
				handler := r.handlers[n.pattern]
				handler(ctx)
			}else {
				ctx.Json(restfulu.NotFound("NOT_FOUND"))
			}
		//}(conn)

	}
}
