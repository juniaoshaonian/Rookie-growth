package webframe

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type router struct {
	trees map[string]*node
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

func newRouter() router {
	return router{
		trees: map[string]*node{},
	}
}

type node struct {
	router     string
	typ        nodeType
	children   map[string]*node
	path       string
	starChild  *node
	regChild   *node
	paramChild *node
	paramName  string
	handler    HanleFunc
	regExpr    *regexp.Regexp
}

func (r *router) addRoute(method string, path string, fn HanleFunc) {
	if path == "" {
		panic("web: 路由是空字符串")
	}
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}
	if path == "/" {
		if root.handler != nil {
			panic("web: 路由冲突[/]")
		}
		root.handler = fn
		return
	}
	if !strings.HasPrefix(path, "/") {
		panic("web: 路由必须以 / 开头")
	}
	if path != "/" && strings.HasSuffix(path, "/") {
		panic("web: 路由不能以 / 结尾")
	}

	segs := strings.Split(strings.Trim(path, "/"), "/")
	for _, seg := range segs {
		if seg == "" {
			panic(fmt.Sprintf("web: 非法路由。不允许使用 //a/b, /a//b 之类的路由, [%s]", path))
		}
		root = root.createOrChild(seg)
	}
	if root.handler != nil {
		panic(fmt.Sprintf("web: 路由冲突[%s]", path))
	}
	root.router = path
	root.handler = fn

}

func (n *node) createOrChild(path string) *node {
	if path == "*" {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.regChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册通配符路由和正则路由 [%s]", path))
		}
		if n.starChild == nil {
			n.starChild = &node{
				typ:  nodeTypeAny,
				path: path,
			}
		}
		return n.starChild
	}
	if strings.HasPrefix(path, ":") && !strings.Contains(path, "(") && !strings.Contains(path, ")") {
		if n.regChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有正则路由。不允许同时注册正则路由和参数路由 [%s]", path))
		}
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和参数路由 [%s]", path))
		}
		if n.paramChild == nil {

			n.paramChild = &node{
				typ:       nodeTypeParam,
				path:      path,
				paramName: path[1:],
			}
		}
		if n.paramChild.path != path {
			panic(fmt.Sprintf("web: 路由冲突，参数路由冲突，已有 %s，新注册 %s", n.paramChild.path, path))
		}
		return n.paramChild
	}

	if strings.HasPrefix(path, ":") && strings.Contains(path, "(") && strings.Contains(path, ")") {
		if n.paramChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有路径参数路由。不允许同时注册正则路由和参数路由 [%s]", path))
		}
		if n.starChild != nil {
			panic(fmt.Sprintf("web: 非法路由，已有通配符路由。不允许同时注册通配符路由和正则路由 [%s]", path))
		}
		if n.regChild == nil {
			start := strings.Index(path, "(")
			key := path[1:start]
			regstr := path[start+1 : len(path)-1]
			reg, err := regexp.Compile(regstr)
			if err != nil {
				log.Fatal(err)
			}
			n.regChild = &node{
				typ:       nodeTypeReg,
				regExpr:   reg,
				paramName: key,
				path:      path,
			}
		}
		return n.regChild
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]
	if !ok {
		child = &node{
			path:     path,
			children: make(map[string]*node),
		}
		n.children[path] = child
	}
	return child
}

func (r *router) findRoute(method, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		return &matchInfo{
			n: root,
		}, true
	}
	segs := strings.Split(strings.Trim(path, "/"), "/")
	mi := &matchInfo{}
	for _, seg := range segs {
		var matchparam bool
		root, matchparam, ok = root.ChildOF(seg)
		if !ok {
			return nil, false
		}
		if matchparam {
			mi.AddValue(root.paramName, seg)
		}
	}
	mi.n = root

	return mi, true

}

func (n *node) ChildOF(path string) (*node, bool, bool) {
	if n.children == nil {
		if n.regChild != nil {
			if n.regChild.regExpr.MatchString(path) {
				return n.regChild, true, true
			}
		}
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		if n.path == "*" {
			return n, false, true
		}

		return n.starChild, false, n.starChild != nil
	}
	res, ok := n.children[path]
	if !ok {
		if n.regChild != nil {
			if n.regChild.regExpr.MatchString(path) {
				return n.regChild, true, true
			}
		}
		if n.paramChild != nil {
			return n.paramChild, true, true
		}
		if n.path == "*" {
			return n, false, true
		}
		return n.starChild, false, n.starChild != nil
	}

	return res, false, ok
}

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) AddValue(key string, value string) {
	if m.pathParams == nil {
		m.pathParams = map[string]string{
			key: value,
		}
	}
	m.pathParams[key] = value
}
