package web

import (
	"fmt"
	"net/http"
	"testing"
)

//Middleware 函数式的责任

type User struct {
	Name string `json:"name"`
}

func TestMiddleware_test(t *testing.T) {

	userMiddle := []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("user的middleware")
				next(ctx)
				ctx.RespData = []byte("xxxxxx")
				fmt.Println("11111111")
			}
		},
	}

	aMiddle := []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				next(ctx)
				ctx.RespData = []byte("nihao")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				next(ctx)
			}
		},
	}
	s := NewHTTPServer()

	s.Get("/user", func(ctx *Context) {

		ctx.RespJSON(200, User{Name: "ser"})
	})
	s.Get("/login", func(ctx *Context) {

		ctx.RespJSON(200, User{Name: "ser"})
	})
	s.Get("/a/*", func(ctx *Context) {

		ctx.RespJSON(200, User{Name: "/a/1"})
	})
	s.UseV1(http.MethodGet, "/user", userMiddle...)

	s.UseV1(http.MethodGet, "/a/*", aMiddle...)

	s.Start(":8081")
}
