package webframe_2

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func Test_router_AddRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodGet,
			path:   "/user/home",
		},
		{
			method: http.MethodGet,
			path:   "/order/detail",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodPost,
			path:   "/login",
		},
		// 通配符测试用例
		{
			method: http.MethodGet,
			path:   "/order/*",
		},
		{
			method: http.MethodGet,
			path:   "/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/*",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc",
		},
		{
			method: http.MethodGet,
			path:   "/*/abc/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},
		// 正则路由
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:name(^.+$)/abc",
		},
	}

	mockfn := func(ctx *Context) {}
	r := newRouter()
	for _, tr := range testRoutes {
		r.addRouter(tr.method, tr.path, mockfn)
	}

	wantRouter := &router{
		trees: map[string]*node{
			http.MethodGet: {
				path: "/",
				children: map[string]*node{
					"user": {
						path: "user",
						children: map[string]*node{
							"home": {path: "home", fn: mockfn, typ: staticNode},
						},
						fn:  mockfn,
						typ: staticNode,
					},
					"order": {
						path: "order",
						children: map[string]*node{
							"detail": {path: "detail", fn: mockfn, typ: staticNode},
						},
						starchild: &node{path: "*", fn: mockfn, typ: starNode},
						typ:       staticNode,
					},
					"param": {
						path: "param",
						paramchild: &node{
							path:  ":id",
							param: "id",
							starchild: &node{
								path: "*",
								fn:   mockfn,
								typ:  starNode,
							},
							children: map[string]*node{"detail": {path: "detail", fn: mockfn, typ: staticNode}},
							fn:       mockfn,
							typ:      paramNode,
						},
					},
				},
				starchild: &node{
					path: "*",
					children: map[string]*node{
						"abc": {
							path:      "abc",
							starchild: &node{path: "*", fn: mockfn, typ: starNode},
							fn:        mockfn,
							typ:       staticNode,
						},
					},
					starchild: &node{path: "*", fn: mockfn, typ: starNode},
					fn:        mockfn,
					typ:       starNode,
				},
				fn:  mockfn,
				typ: staticNode,
			},
			http.MethodPost: {
				path: "/",
				children: map[string]*node{
					"order": {path: "order", children: map[string]*node{
						"create": {path: "create", fn: mockfn, typ: staticNode},
					}},
					"login": {path: "login", fn: mockfn, typ: staticNode},
				},
				typ: staticNode,
			},
			http.MethodDelete: {
				path: "/",
				children: map[string]*node{
					"reg": {
						path: "reg",
						typ:  staticNode,
						regchild: &node{
							path:  ":id(.*)",
							param: "id",
							typ:   regNode,
							fn:    mockfn,
						},
					},
				},
				regchild: &node{
					path:  ":name(^.+$)",
					param: "name",
					typ:   regNode,
					children: map[string]*node{
						"abc": {
							path: "abc",
							fn:   mockfn,
						},
					},
				},
			},
		},
	}
	msg, ok := wantRouter.equal(r)
	assert.True(t, ok, msg)

	// 非法用例
	r = newRouter()

}
func (r router) equal(y router) (string, bool) {
	for k, v := range r.trees {
		yv, ok := y.trees[k]
		if !ok {
			return fmt.Sprintf("目标 router 里面没有方法 %s 的路由树", k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return k + "-" + str, ok
		}
	}
	return "", true
}

func (n *node) equal(y *node) (string, bool) {
	if y == nil {
		return "目标节点为nil", false
	}
	if n.path != y.path {
		return fmt.Sprintf("%s 节点 path 不相等 x %s, y %s", n.path, n.path, y.path), false
	}
	nfn := reflect.ValueOf(n.fn)
	yfn := reflect.ValueOf(y.fn)
	if nfn != yfn {
		return fmt.Sprintf("%s 节点 fn 不相等 x %s, y %s", n.path, nfn.Type().String(), yfn.Type().String()), false
	}
	if len(n.children) != len(y.children) {
		return fmt.Sprintf("%s 子节点长度不等", n.path), false
	}
	if len(n.children) == 0 {
		return "", true
	}
	if n.starchild != nil {
		str, ok := n.starchild.equal(y.starchild)
		if !ok {
			return fmt.Sprintf("%s 通配符节点不匹配 %s", n.path, str), false
		}

	}
	if n.paramchild != nil {
		str, ok := n.paramchild.equal(y.paramchild)
		if !ok {
			return fmt.Sprintf("%s 路径参数节点不匹配 %s", n.path, str), false
		}
	}

	if n.regchild != nil {
		str, ok := n.regchild.equal(y.regchild)
		if !ok {
			return fmt.Sprintf("%s 路径参数节点不匹配 %s", n.path, str), false
		}
	}
	for k, v := range n.children {
		yv, ok := y.children[k]
		if !ok {
			return fmt.Sprintf("%s 目标节点缺少子节点 %s", n.path, k), false
		}
		str, ok := v.equal(yv)
		if !ok {
			return n.path + "-" + str, ok
		}
	}
	return "", true
}

func Test_router_findRoute(t *testing.T) {
	testRoutes := []struct {
		method string
		path   string
	}{
		{
			method: http.MethodGet,
			path:   "/",
		},
		{
			method: http.MethodGet,
			path:   "/user",
		},
		{
			method: http.MethodPost,
			path:   "/order/create",
		},
		{
			method: http.MethodGet,
			path:   "/user/*/home",
		},
		{
			method: http.MethodPost,
			path:   "/order/*",
		},
		// 参数路由
		{
			method: http.MethodGet,
			path:   "/param/:id",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/detail",
		},
		{
			method: http.MethodGet,
			path:   "/param/:id/*",
		},

		// 正则
		{
			method: http.MethodDelete,
			path:   "/reg/:id(.*)",
		},
		{
			method: http.MethodDelete,
			path:   "/:id([0-9]+)/home",
		},
	}

	mockfn := func(ctx *Context) {}

	testCases := []struct {
		name   string
		method string
		path   string
		found  bool
		mi     *matchInfo
	}{
		{
			name:   "method not found",
			method: http.MethodHead,
		},
		{
			name:   "path not found",
			method: http.MethodGet,
			path:   "/abc",
		},
		{
			name:   "root",
			method: http.MethodGet,
			path:   "/",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "/",
					fn:   mockfn,
				},
			},
		},
		{
			name:   "user",
			method: http.MethodGet,
			path:   "/user",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "user",
					fn:   mockfn,
				},
			},
		},
		{
			name:   "no fn",
			method: http.MethodPost,
			path:   "/order",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "order",
				},
			},
		},
		{
			name:   "two layer",
			method: http.MethodPost,
			path:   "/order/create",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "create",
					fn:   mockfn,
				},
			},
		},
		// 通配符匹配
		{
			// 命中/order/*
			name:   "star match",
			method: http.MethodPost,
			path:   "/order/delete",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "*",
					fn:   mockfn,
				},
			},
		},
		{
			// 命中通配符在中间的
			// /user/*/home
			name:   "star in middle",
			method: http.MethodGet,
			path:   "/user/Tom/home",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "home",
					fn:   mockfn,
				},
			},
		},
		{
			// 比 /order/* 多了一段
			name:   "overflow",
			method: http.MethodPost,
			path:   "/order/delete/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "*",
					fn:   mockfn,
				},
			},
		},
		// 参数匹配
		{
			// 命中 /param/:id
			name:   ":id",
			method: http.MethodGet,
			path:   "/param/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: ":id",
					fn:   mockfn,
				},
				params: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /param/:id/*
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/abc",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "*",
					fn:   mockfn,
				},
				params: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /param/:id/detail
			name:   ":id*",
			method: http.MethodGet,
			path:   "/param/123/detail",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: "detail",
					fn:   mockfn,
				},
				params: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /reg/:id(.*)
			name:   ":id(.*)",
			method: http.MethodDelete,
			path:   "/reg/123",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: ":id(.*)",
					fn:   mockfn,
				},
				params: map[string]string{"id": "123"},
			},
		},
		{
			// 命中 /:id([0-9]+)/home
			name:   ":id([0-9]+)",
			method: http.MethodDelete,
			path:   "/123/home",
			found:  true,
			mi: &matchInfo{
				n: &node{
					path: ":id(.*)",
					fn:   mockfn,
				},
				params: map[string]string{"id": "123"},
			},
		},
		{
			// 未命中 /:id([0-9]+)/home
			name:   "not :id([0-9]+)",
			method: http.MethodDelete,
			path:   "/abc/home",
		},
	}

	r := newRouter()
	for _, tr := range testRoutes {
		r.addRouter(tr.method, tr.path, mockfn)
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mi, found := r.findrouter(tc.method, tc.path)
			assert.Equal(t, tc.found, found)
			if !found {
				return
			}
			assert.Equal(t, tc.mi.params, mi.params)
			n := mi.n
			wantVal := reflect.ValueOf(tc.mi.n.fn)
			nVal := reflect.ValueOf(n.fn)
			assert.Equal(t, wantVal, nVal)
		})
	}
}
