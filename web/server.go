package web

import (
	"net"
	"net/http"
)

type HandlerFunc func(ctx Context)

//确保 httpserver 一定实现了 server里面的方法
var _ Server = &HttpServer{}

type Server interface {
	http.Handler
	Start(addr string) error
	// AddRoute 增加路由注册  method方法 path路径 handle业务方法
	addRoute(method, path string, handleFunc HandlerFunc)
	//Get(path string, handleFunc HandlerFunc)
	//Post(path string, handleFunc HandlerFunc)
}

type HttpServer struct {
	*router
}

func NewHTTPServer() *HttpServer {
	return &HttpServer{
		router: newRouter(),
	}
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//你的框架代码在这里
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	h.serve(ctx)
}

func (h *HttpServer) serve(ctx *Context) {
	mi, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n.handler == nil {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		_, _ = ctx.Resp.Write([]byte("not found "))
		return
	}
	ctx.PathParams = mi.pathParams
	mi.n.handler(*ctx)

}
func (h *HttpServer) Start(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return http.Serve(l, h)
}
