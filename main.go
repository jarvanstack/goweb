package main

import (
	"fmt"

	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb6/goweb"
)

func main() {
	w := v1.NewWeb("/bmft")
	w.NewGroup("/v1")
	w.Get("/json", func(ctx *v1.Context) {
		u := &U{}
		ctx.Unmarshal(u)
		ctx.Json(restfulu.Ok(u.Name))
	})
	w.RunHTTP(8001)
}
