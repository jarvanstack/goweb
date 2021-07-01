package main

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb3/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/bmft")
	v1 := web.NewGroup("/v1")
	{
		v2 := v1.NewGroup("/v2")
		{
			v2.Get("/ping", func(ctx *goweb.Context) {
				ctx.Json(restfulu.Ok(ctx.Path))
			})
		}
	}
	web.RunHTTP(8888)
	// http://localhost:8888/bmft/v1/v2/ping
}
