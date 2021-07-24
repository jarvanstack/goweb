package main

import (
	"fmt"
	"strings"

	"github.com/dengjiawen8955/go_utils/restfulu"
	"github.com/dengjiawen8955/goweb/goweb6/goweb"
)

func main() {
	fmt.Printf("%s\n", "hi")
	web := goweb.NewWeb("/bmft")
	v1 := web.NewGroup("/v1")
	{
		v1.Post("/sf", func(ctx *goweb.Context) {
			b := ctx.GetBody()
			fmt.Printf("%s\n", "--------body------")
			fmt.Printf("%s\n", string(b))
			//--处理 1. 拿到boundary
			// buf := bytes.NewBuffer(b)
			// buf.ReadBytes("")
			//返回数据
			// form, err := ctx.GetForm()
			// if err != nil {
			// 	fmt.Printf("%s\n", "bad")
			// } else {
			// 	fmt.Printf("form=%#v\n", form)
			// 	fmt.Printf("files=%#v\n", form.FormFileMap)
			// 	f2 := form.FormFileMap["file2"].Data
			// 	fmt.Printf("f2=%#v\n", f2)
			// 	fmt.Printf("f2=%#v\n", string(f2))
			// 	f3 := form.FormFileMap["img3"]
			// 	file, _ := os.OpenFile(f3.FileName, os.O_RDWR|os.O_CREATE, 0766)
			// 	file.Write(f3.Data)
			// 	file.Close()
			// 	fmt.Printf("datas=%#v\n", form.FormDataMap)
			// 	fmt.Printf("%s\n", "ok")
			// }
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
