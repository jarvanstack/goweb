package main

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/go_utils/testu"
	"github.com/dengjiawen8955/goweb/goweb4/goweb"
	"log"
	"time"
)

func logMiddleware(ctx *goweb.Context)  {
	log.Printf("%s-%s",ctx.Method,ctx.Path)
	start := time.Now().UnixNano()
	ctx.Next()
	log.Printf("---COST=%dms\n", (time.Now().UnixNano()-start)/testu.NS_TO_MS)
}
func main() {
	web := goweb.NewWeb("/bmft")
	v1 := web.NewGroup("/goweb")
	v1.AddMiddleware(logMiddleware)
	{
		v1.Get("/ping", func(ctx *goweb.Context) {
			ctx.Json(restfulu.Ok(ctx.Path))
		})
	}
	web.RunHTTP(8888)
	// http://localhost:8888/bmft/v1/ping
}
