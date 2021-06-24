package goweb

import (
	"net/http"
	"strconv"
)

const (
	BASE_IP = "0.0.0.0:"
	GET     = "GET"
	POST    = "POST"
)

//Web
//Usage:
//w := NewWeb("/bmft")
//w.Get()
//w.Run()
type Web struct {
	r *router
}

func NewWeb(contextPath string) *Web {
	return &Web{r: newRouter(contextPath)}
}
func (w *Web) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	c := newContext(writer, request)
	w.r.handle(c)
}
func (w *Web) Run(port int) error {
	addr := BASE_IP + strconv.Itoa(port)
	//handler HTTP 的关键就是实现 ServeHTTP 这个方法.
	return http.ListenAndServe(addr, w)
}
func (w *Web) Get(path string, chf ContextHandlerFunc) {
	w.r.AddRouter(GET, path, chf)
}
func (w *Web) Post(path string, chf ContextHandlerFunc) {
	w.r.AddRouter(POST, path, chf)
}
