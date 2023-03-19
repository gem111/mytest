package web

import "testing"

func TestLoginPage(t *testing.T) {
	h := NewHTTPServer()
	h.Get("/login", func(ctx *Context) {
		//ctx.Render()
	})
}
