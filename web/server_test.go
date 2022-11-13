package web

import (
	"fmt"
	"testing"
)

func TestServer(t *testing.T) {

	server := web.NewHTTPServer()

	server.Get("/test", func(ctx web.Context) {
		ctx.Resp.Write([]byte("test"))
	})
	server.Get("/param/:id", func(ctx web.Context) {

		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s 参数 %s", ctx.Req.URL.Path, ctx.PathParams)))
	})
	//通配符匹配
	server.Get("/test4/*/*/*", func(ctx web.Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s ", ctx.Req.URL.Path)))
	})
	server.Get("/param2/:id/test", func(ctx web.Context) {

		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s 参数 %s", ctx.Req.URL.Path, ctx.PathParams)))
	})
	server.Get("/param/:id/:username", func(ctx web.Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s 参数 %s", ctx.Req.URL.Path, ctx.PathParams)))
	})
	server.Get("/test1/*/test", func(ctx web.Context) {
		ctx.Resp.Write([]byte("/test/*/test"))
	})
	server.Get("/*", func(ctx web.Context) {
		ctx.Resp.Write([]byte("/*"))
	})
	server.Get("/test2/*/*", func(ctx web.Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s", ctx.Req.URL.Path)))
	})
	server.Get("/test3/*", func(ctx web.Context) {
		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s", ctx.Req.URL.Path)))
	})
	server.Start(":8081")

}
func TestServer2(t *testing.T) {

	server := web.NewHTTPServer()

	//通配符匹配
	server.Get("/test/*/*", func(ctx web.Context) {

		ctx.Resp.Write([]byte(fmt.Sprintf("hello %s 参数 %s", ctx.Req.URL.Path, ctx.PathParams)))
	})
	server.Start(":8081")

}
