package main

import (
	"github.com/dengjiawen8955/goweb/goweb1/goweb"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	web := goweb.NewWeb("/ws")
	web.Get("/ping", func(ctx *goweb.Context) {
		//升级为 websocket
		ws, _ := ctx.NewWs()
		for  {
			msg, _ := ws.ReadMsg()
			ws.WriteMsg(msg)
		}
	})
	web.RunHTTP(8888)
	//ws://localhost:8888/ws/ping
	//http://coolaf.com/tool/chattest
	//测试工具
}
