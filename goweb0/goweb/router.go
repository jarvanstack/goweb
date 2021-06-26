package goweb

import (
	"fmt"
	"log"
)

const (
	keyFmt = "%s-%s"	//比如 get-/user/name
)

type HttpHandler func(ctx *Context)
type router struct {
	httpMap map[string]HttpHandler
	cp      string
	h       *Http
}
func newRouter(contextPath string) *router {
	return &router{httpMap: make(map[string]HttpHandler), cp: contextPath}
}
func (r *router) addRouter(method, urlPath string, httpHandler HttpHandler) {
	key := fmt.Sprintf(keyFmt, method, r.cp+urlPath)
	r.httpMap[key] = httpHandler
}
func (r *router) runHTTP(network,addr string)  {
	log.Printf("network=%s,addr%s\n",network,addr)
	h := newHttp(network, addr,r)
	r.h = h
	h.runHTTP()
}