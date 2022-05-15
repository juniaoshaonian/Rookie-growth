package main

type Routable interface {
	Route(method string,pattern string,handlefn handlerfunc)error
}
type handlerfunc func(c *Context)
type Handler interface {
	ServeHTTP(C *Context)
	Routable
}
