package web

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()

	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个beforme")
				next(ctx)
				fmt.Println("next之后的after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个beforme")
				next(ctx)
				fmt.Println("第二个after")
				ctx.RespData = []byte("111")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {

				fmt.Println("第三个中断了没有")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("我觉得你看不到我")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
