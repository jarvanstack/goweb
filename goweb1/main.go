package main

import (
	"github.com/dengjiawen8955/goweb/goweb1/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/v1")
	web.Get("/ping", func(ctx *goweb.Context) {
		ctx.Json("PONG")
	})
	web.RunHTTP(8888)
	//http://localhost:8888/v1/ping
}