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
	*Group //继承 router 所有方法
	r      *router
	groups []*Group // store all groups
	contextPath string
}
func NewWeb(contextPath string) *Web {
	web := &Web{contextPath: contextPath}
	web.r = newRouter(contextPath,web)
	web.Group = &Group{
		web: web,
	}
	web.groups = []*Group{web.Group}
	return web
}

// NewGroup is defined to create a new Group
// remember all groups share the same Engine instance
func (g *Group) NewGroup(prefix string) *Group {
	web := g.web
	newGroup := &Group{
		prefix: g.prefix + prefix,
		parent: g,
		web:    web,
	}
	web.groups = append(web.groups, newGroup)
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
type Group struct {
	prefix      string        // GET-/bmft 拿些放哪里？
	middlewares []HttpHandler // support middleware
	parent      *Group        // support nesting
	web      *Web             // all groups share a Engine instance
}

//添加中间件
func (g *Group) AddMiddleware(middlewares ...HttpHandler)  {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) addRoute(method string, comp string, handler HttpHandler) {
	pattern :=  g.prefix + comp
	//contextPath 是在 router 里面添加的
	g.web.r.addRouter(method, pattern, handler)
}
func (g *Group) Get(path string,h HttpHandler)  {
	g.addRoute(GET,path,h)
}
func (g *Group) Post(path string,h HttpHandler)  {
	g.addRoute(POST,path,h)
}