# goweb

## goweb0 (http,restful,context)

a simple web framework use tcp

how to use

```go
web := goweb.NewWeb("/v1")
web.Get("/ping", func(ctx *goweb.Context) {
ctx.Json(restfulu.Ok("PONG"))
})
web.RunHTTP(8888)
//http://localhost:8888/v1/ping
```
test

```bash
$ curl localhost:8888/v1/ping
{"code":200,"msg":"OK","data":"PONG"}
```

**performance test**

```go
类型                           RPS
goweb 不解析header:           1917/2088/1836
go 官方 http                  1939/2097/1946
goweb 框架                    1985/1981/2080
```

## goweb1 (websocket)

implements web socket


how to use

```go
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
```

test

open this websocket test website: [http://coolaf.com/tool/chattest](http://coolaf.com/tool/chattest)

input your websocket address 

```bash
ws://localhost:8888/ws/ping
```

result

![web](/img/websocket1.png)

## goweb2 trie router

use trie implements router and params store.

how to use.

```go
web := goweb.NewWeb("/bmft")
web.Get("/:lang/doc", func(ctx *goweb.Context) {
    fmt.Printf("ctx.Params=%#v\n", ctx.Params)
    lang := ctx.Params["lang"]
    fmt.Printf("lang=%#v\n", lang)
    ctx.Json(restfulu.Ok(lang))
})
web.RunHTTP(8888)
```

test 

```bash
$ curl http://localhost:8888/bmft/go/doc
{"code":200,"msg":"OK","data":"doc"}
```

## goweb3 group control


