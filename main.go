package main

import (
	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	w := goweb.NewWeb("/v1")
	w.Post("/json", func(ctx *goweb.Context) {
		u := &User{}
		ctx.UnmarshalJson(u)
		ctx.Json(restfulu.Ok(u.Name))
	})
	w.RunHTTP(8888)
}
