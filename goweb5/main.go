package main

import (
	"fmt"
	"github.com/dengjiawen8955/goweb/goweb5/goweb"
	"time"
)


func main() {
	web := goweb.NewWeb("/bmft")
	v1 := web.NewGroup("/v1")
	{
		v1.Get("/ping", func(ctx *goweb.Context) {
			time.Sleep(time.Second)
			panic(fmt.Errorf("error:%s", "error test"))
		})
	}
	web.RunHTTP(8888)
	// http://localhost:8888/bmft/v1/ping
}
