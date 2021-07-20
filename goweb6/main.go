package main

import (
	"fmt"
	"strings"

	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb6/goweb"
)

func main() {
	web := goweb.NewWeb("/bmft")
	v1 := web.NewGroup("/v1")
	{
		v1.Post("/sf", func(ctx *goweb.Context) {
			b := ctx.Body
			fmt.Printf("%s\n", "--------body------")
			fmt.Printf("%s\n", string(b))
			//--处理 1. 拿到boundary
			// buf := bytes.NewBuffer(b)
			// buf.ReadBytes("")
			//返回数据
			ctx.Json(restfulu.Ok("OK"))
		})
	}
	web.RunHTTP(8888)
	// http://localhost:8888/bmft/v1/sf
}

func getD(ctx *goweb.Context) (string, error) {
	ct, ok := ctx.Headers[goweb.ContentType]
	if !ok {
		return "", fmt.Errorf("ctx.Headers[goweb.ContentType]\n")
	}
	splits := strings.Split(ct, "boundary=")
	if len(splits) != 2 {
		return "", fmt.Errorf("strings.Split(ct, \"boundary=\") not 2\n")
	}
	return splits[1], nil
}
