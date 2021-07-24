package main

import (
	"fmt"

	"github.com/dengjiawen8955/go_utils/restfulu"
	v1 "github.com/dengjiawen8955/goweb/v1"
)

func main() {
	w := v1.NewWeb("/bmft")
	w.NewGroup("/v1")
	w.Get("/ping", func(ctx *v1.Context) {
		cl := ctx.Headers["content-length"]
		fmt.Printf("%s\n", cl)
		ctx.Json(restfulu.Ok("PONG"))
	})
	w.RunHTTP(8001)
}
