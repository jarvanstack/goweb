package main

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb0/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/goweb")
	web.Get("/ping", func(ctx *goweb.Context) {
		ctx.Json(restfulu.Ok("PONG"))
	})
	web.RunHTTP(8888)
	//http://localhost:8888/v1/ping
}