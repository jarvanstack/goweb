package goweb

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restful_util"
	"log"
)

type ContextHandlerFunc func(c *Context)
type router struct {
	handlerFuncMap map[string]ContextHandlerFunc
	contextPath    string
}

//Return a router pointer.
func newRouter(contextPath string) *router {
	return &router{handlerFuncMap: make(map[string]ContextHandlerFunc), contextPath: contextPath}
}

//Add a router.
//if you init contextPath is "/example" engine.AddRouter("GET","/login",handlerFunc)
//you request url should be `ip:port/example/login`
func (this *router) AddRouter(method, urlPath string, contextHandlerFunc ContextHandlerFunc) {
	key := fmt.Sprintf("%s-%s", method, this.contextPath+urlPath)
	this.handlerFuncMap[key] = contextHandlerFunc
}
func (this *router) handle(c *Context) {
	key := fmt.Sprintf("%s-%s", c.Method, c.Path)
	log.Printf("url path=%#v\n", key)
	contextHandlerFunc := this.handlerFuncMap[key]

	//if have this method mapping.
	if contextHandlerFunc != nil {
		contextHandlerFunc(c)
	} else { //if do not have this mapping return 404
		_, _ = c.JSON(restful_util.NotFound(nil))
	}
}
