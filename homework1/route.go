package web

import (
	"fmt"
	"regexp"
	"strings"
)

type router struct {
	// trees 是按照 HTTP 方法来组织的
	// 如 GET => *node
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

// addRoute 注册路由。
// method 是 HTTP 方法
// - 已经注册了的路由，无法被覆盖。例如 /user/home 注册两次，会冲突
// - path 必须以 / 开始并且结尾不能有 /，中间也不允许有连续的 /
// - 不能在同一个位置注册不同的参数路由，例如 /user/:id 和 /user/:name 冲突
// - 不能在同一个位置同时注册通配符路由和参数路由，例如 /user/:id 和 /user/* 冲突
// - 同名路径参数，在路由匹配的时候，值会被覆盖。例如 /user/:id/abc/:id，那么 /user/123/abc/456 最终 id = 456
func (r *router) addRoute(method string, path string, handler HandleFunc) {
	//TODO implement me

	if path == "" {
		panic("路径不能为空")
	}

	if path[0] != '/' && path != "/" {
		panic("路径必须[/]开头")
	}
	if path[len(path)-1] == '/' && path != "/" {
		panic("路径不能以[/]结尾")
	}
	root, ok := r.trees[method]
	if !ok {
		//说明还没有建立根结点  先创建[get]跟节点
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path == "/" {
		if root.handler != nil {
			panic("不允许重复注册[/]路径")
		}
		root.handler = handler
		return
	}

	//切割这个path  在get数组下面创建 路径树
	segs := strings.Split(path[1:], "/")

	for _, seg := range segs {
		if seg == "" {
			panic("不能注册连续的//")
		}
		children := root.childOrCreate(seg)
		root = children
		//print("执行root等于", root)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("路由[%s]已被注册", path))
	}
	root.handler = handler
}

// findRoute 查找对应的节点
// 注意，返回的 node 内部 HandleFunc 不为 nil 才算是注册了路由
func (r *router) findRoute(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}

	if path == "/" {
		return &matchInfo{n: root}, true
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	for _, seg := range segs {
		var child *node
		child, ok := root.childOf(seg)
		if !ok {
			if root.typ == nodeTypeAny {
				return &matchInfo{
					n: root,
				}, true
			}
			return nil, false
		}
		//if root.typ == nodeTypeParam {
		//	//赋值
		//	if pathParams == nil {
		//		pathParams = make(map[string]string)
		//	}
		//	pathParams[child.path[1:]] = seg
		//}
		if child.paramName != "" {
			mi.addValue(child.paramName, seg)
		}
		root = child
	}
	mi.n = root
	return mi, true
}

type nodeType int

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

// node 代表路由树的节点
// 路由树的匹配顺序是：
// 1. 静态完全匹配
// 2. 正则匹配，形式 :param_name(reg_expr)
// 3. 路径参数匹配：形式 :param_name
// 4. 通配符匹配：*
// 这是不回溯匹配
type node struct {
	typ nodeType

	path string
	// children 子节点
	// 子节点的 path => node
	children map[string]*node
	// handler 命中路由之后执行的逻辑
	handler HandleFunc

	// 通配符 * 表达的节点，任意匹配
	starChild *node

	paramChild *node
	// 正则路由和参数路由都会使用这个字段
	paramName string

	// 正则表达式
	regChild *node
	regExpr  *regexp.Regexp
}

// child 返回子节点
// 第一个返回值 *node 是命中的节点
// 第二个返回值 bool 代表是否命中
func (n *node) childOf(path string) (*node, bool) {
	if n.children == nil {
		return n.childOfNonStatic(path)
	}
	res, ok := n.children[path]
	if !ok {
		return n.childOfNonStatic(path)
	}
	return res, ok
}

func (n *node) childOfNonStatic(path string) (*node, bool) {

	if n.regChild != nil {
		if n.regChild.regExpr.Match([]byte(path)) {
			return n.regChild, true
		}
	}
	if n.paramChild != nil {
		return n.paramChild, true
	}
	return n.starChild, n.starChild != nil
}

// childOrCreate 查找子节点，
// 首先会判断 path 是不是通配符路径
// 其次判断 path 是不是参数路径，即以 : 开头的路径
// 最后会从 children 里面查找，
// 如果没有找到，那么会创建一个新的节点，并且保存在 node 里面
func (n *node) childOrCreate(seg string) *node {

	if seg[0] == ':' {
		paramName, expr, isReg := n.parseParam(seg)
		if isReg {
			//正则匹配
			return n.childOrCreateReg(paramName, expr, paramName)

		}
		return n.childOrCreateParam(seg, paramName)
	}
	if seg == "*" {
		if n.paramChild != nil {
			panic("不可以同时注册通配符路由和参数路由")
		}
		if n.regChild != nil {
			panic(fmt.Sprintf("不可以同时注册通配符和正则路由"))
		}
		if n.starChild == nil {
			n.starChild = &node{
				typ:  nodeTypeAny,
				path: seg,
			}
		}
		return n.starChild
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
func (n *node) childOrCreateReg(path, expr, paramName string) *node {
	if n.paramChild != nil {
		panic(fmt.Sprintf("有参数路由了，不可以同时注册正则和参数路由 %s", path))
	}
	if n.starChild != nil {
		panic(fmt.Sprintf("已有通配符路由，不可以同时注册通配符和参数路由 %s", path))
	}

	if n.regChild != nil {
		if n.regChild.regExpr.String() != expr || n.paramName != paramName {
			panic(fmt.Sprintf("已有参数路由 %s ， 新注册 %s", n.paramChild.path, path))
		}
	} else {
		regExpr, err := regexp.Compile(expr)
		if err != nil {
			panic(fmt.Errorf("正则表达式错误 %w", err))
		}
		n.regChild = &node{
			typ:       nodeTypeReg,
			path:      path,
			paramName: paramName,
			regExpr:   regExpr,
		}
	}
	return n.regChild

}
func (n *node) childOrCreateParam(path, paramName string) *node {
	if n.regChild != nil {
		panic(fmt.Sprintf("有正则路由了，不可以同时注册正则和参数路由 %s", path))
	}
	if n.starChild != nil {
		panic(fmt.Sprintf("已有通配符路由，不可以同时注册通配符和参数路由 %s", path))
	}

	if n.paramChild != nil {
		if n.paramChild.path != path {
			panic(fmt.Sprintf("已有参数路由 %s ， 新注册 %s", n.paramChild.path, path))
		}
	} else {
		n.paramChild = &node{
			typ:       nodeTypeParam,
			path:      path,
			paramName: paramName,
		}
	}
	return n.paramChild
}

func (n *node) parseParam(path string) (string, string, bool) {
	// 去除 :
	path = path[1:]
	segs := strings.SplitN(path, "(", 2)
	if len(segs) == 2 {
		expr := segs[1]
		if strings.HasSuffix(expr, ")") {
			return segs[0], expr[:len(expr)-1], true
		}
	}
	return path, "", false
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key, val string) {
	if m.pathParams == nil {
		m.pathParams = map[string]string{}
	}
	m.pathParams[key] = val
}
