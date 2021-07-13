package v1

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/go_utils/testu"
	"github.com/dengjiawen8955/go_utils/throwu"
	"log"
	"strconv"
	"time"
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
	*Group           //继承 router 所有方法
	r               *router
	groups          []*Group // store all groups
	contextPath     string
	IsLogMiddleware bool //是否日志默认打开
}
func NewWeb(contextPath string) *Web {
	web := &Web{
		contextPath:     contextPath,
		IsLogMiddleware: true,
	}
	web.r = newRouter(contextPath,web)
	web.Group = &Group{
		web: web,
	}
	web.groups = []*Group{web.Group}
	//添加错误处理
	web.AddMiddleware(errMiddleware)
	//添加日志打印
	if web.IsLogMiddleware {
		web.AddMiddleware(logMiddleware)
	}
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
//log记录每次请求信息中间件
func logMiddleware(ctx *Context)  {
	start := time.Now().UnixNano()
	ctx.Next()
	cur := time.Now().UnixNano()
	log.Printf("%s|COST=%dms|%s|%s",
		time.Now().Format(testu.TIME_FOMART),
		(cur-start)/testu.NS_TO_MS,
		ctx.Method,
		ctx.Path,
		)
}
//错误处理中间件,所有的错误都可以通过这个中间件处理
func errMiddleware(ctx *Context)  {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%s\n\n", throwu.Trace(err))
			ctx.Json(restfulu.ServerError("SERVER_ERROR"))
		}
	}()
	ctx.Next()
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
