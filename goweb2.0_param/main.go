package main

import (
	"fmt"
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb3/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/bmft")
	web.Get("/:lang/doc", func(ctx *goweb.Context) {
		fmt.Printf("ctx.Params=%#v\n", ctx.Params)
		lang := ctx.Params["lang"]
		fmt.Printf("lang=%#v\n", lang)
		ctx.Json(restfulu.Ok(lang))
	})
	web.RunHTTP(8888)
}
