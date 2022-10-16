package webframe_2

import (
	"net"
	"net/http"
)

type HandleFunc func(c *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	AddRouter(method string, path string, fn HandleFunc) error
}

type HttpServer struct {
	r    router
	mdls []Middleware
}

func (h *HttpServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}
	root := h.Server
	for i := len(h.mdls) - 1; i >= 0; i++ {
		root = h.mdls[i](root)
	}
	root(ctx)
}

func (h *HttpServer) Server(ctx *Context) {
	m, ok := h.r.findrouter(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok {
		ctx.Resp.WriteHeader(http.StatusNotFound)
		ctx.Resp.Write([]byte("router not found"))
	}
	ctx.paramMap = m.params
	m.n.fn(ctx)
}

func (h *HttpServer) GET(Path string, fn HandleFunc) error {
	return h.AddRouter(http.MethodGet, Path, fn)
}
func (h *HttpServer) POST(Path string, fn HandleFunc) error {
	return h.AddRouter(http.MethodPost, Path, fn)
}

func (h *HttpServer) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	http.Serve(listener, h)
}

func (h *HttpServer) AddRouter(method string, path string, fn HandleFunc) error {
	return h.r.addRouter(method, path, fn)
}
