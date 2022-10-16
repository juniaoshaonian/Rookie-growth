package webframe_2

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"webframe/errs"
)

// 本次路由树支持，通配符，路由参数，正则路由
type router struct {
	trees map[string]*node
}

func newRouter() router {
	return router{
		trees: make(map[string]*node),
	}
}

type nodeTyp int

const (
	staticNode = iota
	regNode
	paramNode
	starNode
)

type node struct {
	typ        nodeTyp
	path       string
	children   map[string]*node
	starchild  *node
	paramchild *node
	regchild   *node
	reg        *regexp.Regexp
	param      string
	fn         HandleFunc
}

// 添加路由

func (r *router) addRouter(method string, path string, fn HandleFunc) error {
	root, ok := r.trees[method]
	if !ok {
		root = &node{
			path: "/",
		}
		r.trees[method] = root
	}

	if path == "/" {
		if root.fn != nil {
			return errs.NewErrPathCoverage(path)
		}
		root.fn = fn
		return nil
	}
	if !strings.HasPrefix(path, "/") {
		return errors.New(fmt.Sprintf("未以/ 开头，%v", path))
	}
	if strings.HasSuffix(path, "/") {
		return errors.New(fmt.Sprintf("不能以/为结尾 %v", path))
	}
	path = strings.Trim(path, "/")
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		if seg == "" {
			return errors.New("path中间不能有重复的/")
		}
		var err error
		root, err = root.CreateOrChild(seg)
		if err != nil {
			return err
		}
	}
	if root.fn != nil {
		return errs.NewErrPathCoverage(path)
	}
	root.fn = fn
	return nil
}

func (n *node) CreateOrChild(path string) (*node, error) {
	if path == "*" {
		if n.regchild != nil {
			return nil, errs.NewErrNodeConflict("正则", "通配符")
		}
		if n.paramchild != nil {
			return nil, errs.NewErrNodeConflict("路由参数", "通配符")
		}
		if n.starchild == nil {
			new_child := &node{
				typ:  starNode,
				path: path,
			}
			n.starchild = new_child
		}
		return n.starchild, nil
	}
	if strings.HasPrefix(path, ":") && !strings.Contains(path, "(") {
		if n.regchild != nil {
			return nil, errs.NewErrNodeConflict("正则", "路由参数")
		}
		if n.starchild != nil {
			return nil, errs.NewErrNodeConflict("路由参数", "通配符")
		}
		if n.paramchild == nil {
			new_node := &node{
				typ:   paramNode,
				path:  path,
				param: path[1:],
			}

			n.paramchild = new_node
		}
		return n.paramchild, nil
	}
	if strings.HasPrefix(path, ":") && strings.Contains(path, "(") {
		if n.starchild != nil {
			return nil, errs.NewErrNodeConflict("统配符", "正则")
		}
		if n.paramchild != nil {
			return nil, errs.NewErrNodeConflict("路由参数", "正则")
		}
		if n.regchild == nil {
			start := 1
			for start < len(path) && path[start] != '(' {
				start++
			}
			end := start
			for end < len(path) && path[end] != ')' {
				end++
			}
			param := path[1:start]
			reg, err := regexp.Compile(path[start+1 : end])
			if err != nil {
				return nil, err
			}
			new_node := &node{
				typ:   regNode,
				reg:   reg,
				param: param,
				path:  path,
			}
			n.regchild = new_node
		}
		return n.regchild, nil
	}
	if n.children == nil {
		n.children = make(map[string]*node)
	}
	child, ok := n.children[path]

	if !ok {
		child = &node{
			typ:  staticNode,
			path: path,
		}
		n.children[path] = child
	}
	return child, nil
}

func (r *router) findrouter(method string, path string) (*matchInfo, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	if path == "/" {
		if root.fn != nil {
			return &matchInfo{
				n: root,
			}, true
		} else {
			return nil, false
		}
	}
	path = strings.Trim(path, "/")
	sets := strings.Split(path, "/")
	m := &matchInfo{}
	for _, seg := range sets {
		mismatched, param, n := root.childOf(seg)
		if !mismatched {
			return nil, false
		}
		if param {
			m.Set(n.param, seg)
		}
		root = n
	}
	m.n = root
	return m, true
}

func (n *node) childOf(path string) (bool, bool, *node) {
	if n.children != nil {
		child, ok := n.children[path]
		if ok {
			return true, false, child
		}
	}
	if n.regchild != nil {
		ok := n.regchild.reg.MatchString(path)
		if ok {
			return true, true, n.regchild
		}
	}
	if n.paramchild != nil {
		return true, true, n.paramchild
	}
	if n.starchild != nil {
		return true, false, n.starchild
	}
	if n.path == "*" {
		return true, false, n
	}
	return false, false, nil
}

type matchInfo struct {
	n      *node
	params map[string]string
}

func (m *matchInfo) Set(key string, value string) {
	if m.params == nil {
		m.params = make(map[string]string)
	}
	m.params[key] = value
}
