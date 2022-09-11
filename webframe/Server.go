package webframe

import "net/http"

type HanleFunc func(c *Context)

type Server interface {
	http.Handler
	Start(addr string) error
	AddRoute(method string, path string)
}
