package goweb

import "net/http"

//Context type
type Context struct {
	//origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	//request info
	Path       string
	Method     string
	StatusCode int
}

//Return context.
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Req:        r,
		Writer:     w,
		Path:       r.URL.Path,
		Method:     r.Method,
		StatusCode: 200,
	}
}

//Send Json
func (c *Context) JSON(jsonBytes []byte) (int, error) {
	c.Writer.Header().Set("Content-Type", "application/json")
	n, err := c.Writer.Write(jsonBytes)
	return n, err
}

//Send text
func (c *Context) Text(jsonBytes []byte) (int, error) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	n, err := c.Writer.Write(jsonBytes)
	return n, err
}

//Get a form value by key
func (c *Context) GetFrom(key string) string {
	return c.Req.Form.Get(key)
}

//Get a url query value by key
func (c *Context) GetQuery(key string) string {
	return c.Req.URL.Query().Get(key)
}

//Set response header status code.
//【原来浏览器自带的code，这个应该是顶层的 code,这玩意在 header 里面】
func (c *Context) SetStatus(statusCode int) {
	c.StatusCode = statusCode
	c.Writer.WriteHeader(statusCode)
}
