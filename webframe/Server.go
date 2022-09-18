package webframe

import (
	"net"
	"net/http"
)

type HanleFunc func(c *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	addRoute(method string, path string)
}

type HttpServr struct {
	router
	mdl []Middleware
}

func (h *HttpServr) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(listener, h)
}

func (h *HttpServr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Req:  r,
		Resp: w,
	}
	root := h.Serve
	for i := len(h.mdl) - 1; i >= 0; i-- {
		root = h.mdl[i](root)
	}
	root(ctx)

}

func (h *HttpServr) Serve(ctx *Context) {
	mi, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n.handler == nil {
		ctx.ResponseCode = 404
		ctx.RespsonseDate = []byte("not found")
		return
	}

	mi.n.handler(ctx)
	ctx.MatchRouter = mi.n.router

}
