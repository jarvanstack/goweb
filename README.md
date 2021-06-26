# goweb

## goweb0 (http,restful,context)

a simple web framework use tcp

how to use

```go
web := goweb.NewWeb("/v1")
web.Get("/ping", func(ctx *goweb.Context) {
    ctx.Json("PONG")
})
web.RunHTTP(8888)
//http://localhost:8888/v1/ping
```
test

```bash
$ curl localhost:8888/v1/ping
{"code":200,"msg":"OK","data":"PONG"}
```

## goweb1 (websocket)

implements web socket

