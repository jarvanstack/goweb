package main

import (
	"github.com/dengjiawen8955/go_utils/restful_util"
	"github.com/dengjiawen8955/goweb/goweb0/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/v1")
	web.Get("/ping", func(ctx *goweb.Context) {
		ctx.Json(restful_util.Ok("PONG"))
	})
	web.RunHTTP(8888)
	//http://localhost:8888/v1/ping
}