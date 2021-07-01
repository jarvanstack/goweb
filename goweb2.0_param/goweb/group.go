package goweb


type Group struct {
	prefix string
	middlewares []HttpHandler
	parent *Group
	web *Web
}
