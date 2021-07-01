package goweb

import (
	"strconv"
)

const (
	GET  = "GET"
	POST = "POST"
)
var (
	network = "tcp"
	baseIp = "0.0.0.0:"
)

type Web struct {
	*RouterGroup //继承 router 所有方法
	r *router
	groups []*RouterGroup // store all groups
}
func NewWeb(contextPath string) *Web {
	web := &Web{
		r: newRouter(contextPath),
	}
	web.RouterGroup = &RouterGroup{
		web: web,
	}
	web.groups = []*RouterGroup{web.RouterGroup}
	return web
}

// NewGroup is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := group.web
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		web: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}


func (w *Web)Get(path string ,h HttpHandler)  {
	w.r.addRouter(GET,path,h)
}
func (w *Web)Post(path string ,h HttpHandler)  {
	w.r.addRouter(POST,path,h)
}
func (w *Web)RunHTTP(port int)  {
	w.r.runHTTP(network,baseIp+strconv.Itoa(port))
}
//----------
type RouterGroup struct {
	prefix      string
	middlewares []HttpHandler // support middleware
	parent      *RouterGroup  // support nesting
	web      *Web       // all groups share a Engine instance
}
func (g *RouterGroup) addRoute(method string, comp string, handler HttpHandler) {
	pattern :=  g.prefix + comp
	//contextPath 是在 router 里面添加的
	g.web.r.addRouter(method, pattern, handler)
}
func (g *RouterGroup) Get(path string,h HttpHandler)  {
	g.addRoute(GET,path,h)
}
func (g *RouterGroup) Post(path string,h HttpHandler)  {
	g.addRoute(POST,path,h)
}