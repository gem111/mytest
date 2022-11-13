package web

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net/http"
	"strings"
)

//用来支持对路由书的操作
//代表路由书  森林
type router struct {
	//beego gin method  对应一棵树
	//get 有一棵树 post也有一棵树
	//http method =》路由书根结点
	trees map[string]*node
}

//提供一个注册创建的方法
func newRouter() *router {

	return &router{
		trees: map[string]*node{},
	}

}

func (r *router) Post(path string, handle HandlerFunc) {
	r.addRoute(http.MethodPost, path, handle)
}
func (r *router) Get(path string, handle HandlerFunc) {
	r.addRoute(http.MethodGet, path, handle)
}
func (r *router) addRoute(method, path string, handle HandlerFunc) {
	//TODO implement me
	root, ok := r.trees[method]
	if !ok {
		//说明还没有建立根结点  先创建[get]跟节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path == "" {
		panic("路径不能为空")
	}

	if path[0] != '/' && path != "/" {
		panic("路径必须[/]开头")
	}
	if path[len(path)-1] == '/' && path != "/" {
		panic("路径不能以[/]结尾")
	}

	if path == "/" {
		if root.handler != nil {
			panic("不允许重复注册[/]路径")
		}
		root.handler = handle
		return
	}

	//切割这个path  在get数组下面创建 路径树
	segs := strings.Split(path[1:], "/")

	for _, seg := range segs {
		if seg == "" {
			panic("不能注册连续的//")
		}
		children := root.chilOrCreate(seg)
		root = children
		//print("执行root等于", root)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("路由[%s]已被注册", path))
	}
	root.handler = handle
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	pathParams := make(map[string]string)
	for _, seg := range segs {
		var child *node
		child, paramChild, ok := root.childOf(seg)
		fmt.Println(seg, "---->", child, paramChild, ok)
		if !ok {
			return nil, false
		}

		//等于通配符直接返回
		if root.typ == nodeTypeAny {
			spew.Println("nihao ", child.startChild)
			return &matchInfo{
				n: root,
			}, true
		}
		if paramChild {
			fmt.Println(seg)
			//赋值
			if pathParams == nil {
				pathParams = make(map[string]string)
			}
			pathParams[child.path[1:]] = seg

		}

		root = child
	}
	return &matchInfo{
		n:          root,
		pathParams: pathParams,
	}, true
}

func (n *node) childOf(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}

	res, ok := n.children[path]
	if !ok {
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		return n.startChild, false, n.startChild != nil
	}
	return res, false, ok
}
func (n *node) chilOrCreate(seg string) *node {
	if seg[0] == ':' {
		if n.startChild != nil {
			panic("不可以同时注册通配符路由和参数路由")
		}
		n.paramChild = &node{
			typ:  nodeTypeParam,
			path: seg,
		}
		return n.paramChild
	}
	if seg == "*" {
		if n.paramChild != nil {
			panic("不可以同时注册通配符路由和参数路由")
		}
		n.startChild = &node{
			typ:  nodeTypeAny,
			path: seg,
		}
		return n.startChild
	}
	if n.children == nil {
		//第一次进来 没有创建 就创建一下
		n.children = make(map[string]*node)
	}
	res, ok := n.children[seg]
	//没有children[user]这一层 就创建一层
	if !ok {
		res = &node{
			typ:  nodeTypeStatic,
			path: seg,
		}
		n.children[seg] = res
	}
	//有就直接返回
	return res
}

const (
	// 静态路由
	nodeTypeStatic = iota
	// 正则路由
	nodeTypeReg
	// 路径参数路由
	nodeTypeParam
	// 通配符路由
	nodeTypeAny
)

type node struct {
	path string

	typ int
	//子 path到子节点的映射
	children map[string]*node

	//通配符
	startChild *node
	//参数匹配
	paramChild *node

	//缺一个代表路由注册的逻辑
	handler HandlerFunc
}
