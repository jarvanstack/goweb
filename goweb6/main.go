package main

import (
	"fmt"
	"strings"

	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb6/goweb"
)

type U struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	fmt.Printf("%s\n", "hi")
	web := goweb.NewWeb("/bmft")
	v1 := web.NewGroup("/v1")
	{
		v1.Post("/json", func(ctx *goweb.Context) {
			b := ctx.GetBody()
			fmt.Printf("%s\n", "--------body------")

			fmt.Printf("%s\n", string(b))
			u := &U{}
			ctx.Unmarshal(u)
			ctx.Json(restfulu.Ok(u.Name))

		})
		v1.Post("/fd", func(ctx *goweb.Context) {
			fmt.Printf("%s\n", "--------body------")
			b := ctx.GetBody()
			fmt.Printf("%s\n", string(b))
			f, err := ctx.GetForm()
			fmt.Printf("err: %v\n", err)
			fd := f.FormDataMap["fd1"]
			fmt.Printf("fd: %v\n", fd)
			fd2 := f.FormDataMap["fd2"]
			fmt.Printf("fd2: %v\n", fd2)
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
