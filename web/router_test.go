package web

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestRouter_AddRoute(t *testing.T) {
	//testRoutes := []struct {
	//	method string
	//	path   string
	//}{
	//	{
	//		method: http.MethodGet,
	//		path:   "/user/home",
	//	},
	//}

	//等于这样
	type testRouterTmp []struct {
		method string
		path   string
	}
	testRoutes := testRouterTmp{
		{
			path:   "/",
			method: http.MethodGet,
		},
		{
			path:   "/user",
			method: http.MethodGet,
		},
		//{
		//	path:   "",
		//	method: http.MethodGet,
		//},
		//{
		//	path:   "/user/////aaa",
		//	method: http.MethodGet,
		//},
		//{
		//	path:   "/user/",
		//	method: http.MethodGet,
		//},
		//{
		//	path:   "aaaa/user/",
		//	method: http.MethodGet,
		//},
		//{
		//	path:   "/user/post",
		//	method: http.MethodPost,
		//},
		{
			path:   "/user/home",
			method: http.MethodGet,
		},
		//{
		//	path:   "/del",
		//	method: http.MethodDelete,
		//},
	}

	r := newRouter()

	mockHandler := func(ctx Context) {}
	//注册路由
	for _, router := range testRoutes {
		r.Get(router.path, mockHandler)
	}

	//模拟出来的路由树的样子，看路由树实际生成的和我们的是否一致
	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: &node{
				path:    "/",
				handler: mockHandler,
				children: map[string]*node{
					"user": &node{
						path:    "user",
						handler: mockHandler,
						children: map[string]*node{
							"home": &node{
								path:    "home",
								handler: mockHandler,
							},
						},
					},
				},
			},
			//http.MethodPost: {
			//	path: "/",
			//	children: map[string]*node{
			//		"user": &node{
			//			path: "user",
			//			children: map[string]*node{
			//				"post": &node{
			//					path:    "post",
			//					handler: mockHandler,
			//				},
			//			},
			//		},
			//	},
			//},
			//http.MethodDelete: {
			//	path: "/",
			//	children: map[string]*node{
			//		"del": &node{
			//			path:    "del",
			//			handler: mockHandler,
			//		},
			//	},
			//},
		},
	}
	msg, ok := wantRouter.equal(r)
	//在这里断言两者相等
	assert.True(t, ok, msg)

	r = newRouter()
	assert.Panicsf(t, func() {
		r.Get("aaaa/user", mockHandler)
	}, "路径必须[/]开头")
	assert.Panicsf(t, func() {
		r.Get("/aaaa/user/", mockHandler)
	}, "路径不能以[/]结尾")
	assert.Panicsf(t, func() {
		r.Get("/aaaa/u//ser", mockHandler)
	}, "不能注册连续的//")
	assert.Panicsf(t, func() {
		r.Get("", mockHandler)
	}, "路径不能为空")
	assert.Panicsf(t, func() {
		r.Get("/", mockHandler)
		r.Get("/", mockHandler)
	}, "不允许重复注册[/]路径")
	assert.Panicsf(t, func() {
		r.Get("/aaa", mockHandler)
		r.Get("/aaa", mockHandler)
	}, "重复注册路由")
	assert.Panicsf(t, func() {
		r.Get("/param/:id", mockHandler)
		r.Get("/param/*", mockHandler)
	}, "路由注册冲突")

}

func (r *router) equal(y *router) (string, bool) {
	for k, v := range r.trees {
		dst, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("找不到对应的 http  method"), false
		}
		msg, equal := v.equal(dst)
		//v,dst 要想等
		if !equal {
			return msg, false
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if n.path != y.path {
		return fmt.Sprintf("节点路径不匹配"), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("子节点数量不相同"), false
	}
	//比对handle
	nhanlder := reflect.ValueOf(n.handler)
	yhanlder := reflect.ValueOf(y.handler)
	if nhanlder != yhanlder {
		return fmt.Sprintf("handler 不想等"), false
	}
	for path, c := range n.children {
		dst, ok := y.children[path]
		if !ok {
			return fmt.Sprintf("子节点 %s 不存在", path), false
		}

		msg, ok := c.equal(dst)
		if !ok {
			return msg, false
		}
	}

	if n.startChild != nil {
		str, ok := n.startChild.equal(y.startChild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.path, str), false
		}
	}

	return "", true
}

func TestRouter_findRouteWildcard(t *testing.T) {

	testRoute := []struct {
		method string
		path   string
	}{
		//注册多个通配符
		{
			method: http.MethodGet,
			path:   "/test/*/*",
		},
	}
	r := newRouter()

	mockHandler := func(ctx Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name     string
		method   string
		path     string
		wanFound bool //有没有找到
		wanNode  *node
	}{
		{
			name:     "/* 多个",
			method:   http.MethodGet,
			path:     "/",
			wanFound: true,
			wanNode: &node{
				path: "/",
				children: map[string]*node{
					"test": &node{
						path: "test",
						startChild: &node{
							path: "*",
							startChild: &node{
								path:    "*",
								handler: mockHandler,
							},
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wanFound, found)
			if !found {
				return
			}
			msg, ok := tc.wanNode.equal(n.n)
			assert.True(t, ok, msg)
		})
	}
}

func TestRouter_findRoute(t *testing.T) {

	testRoute := []struct {
		method string
		path   string
	}{
		//{
		//	method: http.MethodDelete,
		//	path:   "/",
		//},
		//{
		//	method: http.MethodGet,
		//	path:   "/order/detail",
		//},
		//注册多个通配符
		{
			method: http.MethodGet,
			path:   "/test/*/*",
		},
	}
	r := newRouter()

	mockHandler := func(ctx Context) {}
	for _, route := range testRoute {
		r.addRoute(route.method, route.path, mockHandler)
	}

	testCases := []struct {
		name     string
		method   string
		path     string
		wanFound bool //有没有找到
		wanNode  *node
	}{
		{
			name:     "/* 多个",
			method:   http.MethodGet,
			path:     "/",
			wanFound: true,
			wanNode: &node{
				path: "/",
				startChild: &node{
					path:    "*",
					handler: mockHandler,
					startChild: &node{
						path:    "*",
						handler: mockHandler,
						startChild: &node{
							path:    "test",
							handler: mockHandler,
						},
					},
				},
				handler: mockHandler,
			},
		},
		//没有对应的method
		{
			name:   "method not found ",
			method: http.MethodHead,
			path:   "/order/detail",
		},
		//{
		//	//命中了 但是没有handler
		//	name:     "order detail",
		//	method:   http.MethodGet,
		//	path:     "/order",
		//	wanFound: true,
		//	wanNode: &node{
		//		//handler: mockHandler,
		//		path: "order",
		//		children: map[string]*node{
		//			"detail": &node{
		//				handler: mockHandler,
		//				path:    "detail",
		//			},
		//		},
		//	},
		//},
		{
			//注册一个根结点
			name:     "root",
			method:   http.MethodDelete,
			path:     "/",
			wanFound: true,
			wanNode: &node{
				handler: mockHandler,
				path:    "/",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, found := r.findRoute(tc.method, tc.path)
			assert.Equal(t, tc.wanFound, found)
			if !found {
				return
			}
			msg, ok := tc.wanNode.equal(n.n)
			assert.True(t, ok, msg)
			//assert.Equal(t, tc.wanNode.path, n.n.path)
			//nhanlder := reflect.ValueOf(n.n.handler)
			//yhanlder := reflect.ValueOf(tc.wanNode.handler)
			//assert.True(t, yhanlder == nhanlder)
		})
	}
}
