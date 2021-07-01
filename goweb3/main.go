package main

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb2/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/bmft")
	web.Get("/v1/doc", func(ctx *goweb.Context) {
		ctx.Json(restfulu.Ok(ctx.Path))
	})
	web.RunHTTP(8888)
	// http://localhost:8888/bmft/v1/doc
}
