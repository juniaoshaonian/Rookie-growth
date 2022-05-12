package main

type Handler interface {
	ServeHTTP(ctx *Context)
	Routable
}
type handlerfunc func(ctx *Context)
