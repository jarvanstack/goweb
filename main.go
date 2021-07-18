package main

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb"
)

func main() {
	w := goweb.NewWeb("/v1")
	w.Get("/ping", func(ctx *goweb.Context) {
		ctx.Json(restfulu.Ok("PONG"))
	})
	w.RunHTTP(8888)
}
