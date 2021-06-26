package goweb

import "strconv"

const (
	get     = "GET"
	post    = "POST"
)
var (
	network = "tcp"
	baseIp = "0.0.0.0:"
)

type Web struct {
	r *router
}
func NewWeb(contextPath string) *Web {
	return &Web{r: newRouter(contextPath)}
}
func (w *Web)Get(path string ,h HttpHandler)  {
	w.r.addRouter(get,path,h)
}
func (w *Web)Post(path string ,h HttpHandler)  {
	w.r.addRouter(post,path,h)
}
func (w *Web)RunHTTP(port int)  {
	w.r.runHTTP(network,baseIp+strconv.Itoa(port))
}