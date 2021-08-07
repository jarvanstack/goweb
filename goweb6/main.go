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
		v1.Get("/ping", func(ctx *goweb.Context) {
			ctx.Json(restfulu.Ok("PONG"))
		})
		v1.Post("/json", func(ctx *goweb.Context) {
			b := ctx.GetBody()
			fmt.Printf("%s\n", "--------body------")

			fmt.Printf("%s\n", string(b))
			u := &U{}
			ctx.UnmarshalJson(u)
			ctx.Json(restfulu.Ok(u.Name))

		})
		v1.Post("/sf", func(ctx *goweb.Context) {
			fmt.Printf("%s\n", "--------body------")
			b := ctx.GetBody()
			fmt.Printf("%s\n", string(b))
			f, _ := ctx.GetForm()
			f2, _ := f.GetFile("f2")
			// f2f, _ := os.Create("c/" + f2.FileName)
			// f2f.Write(f2.Data)
			// f2f.Close()
			// fmt.Printf("err: %v\n", err)
			// ff, err2 := f.GetFile("f1")
			// fmt.Printf("err2: %v\n", err2)
			// fmt.Println(ff.FileName)
			// ff, _ := f.GetFile("f2")
			// fmt.Printf("err2: %v\n", err2)
			// fmt.Println(ff.FileName)
			fmt.Printf("%s\n", "--------p--------")
			fmt.Print(string(f2.Data))
			fmt.Printf("%s\n", "--------p--------")

			ctx.Json(restfulu.Ok(f2))
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
