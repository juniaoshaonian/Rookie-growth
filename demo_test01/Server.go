package main

import "net/http"

type Routable interface {
	Route(method string,pattern string,handlerfn handlerfunc )
}
type Server interface {
	Routable
	Start(address string)error
}
type sdkhttpServer struct {
	name string
	handler Handler
	root filter
}
func (s *sdkhttpServer)Route(method string,pattern string,handlerfn handlerfunc){
	s.handler.Route(method,pattern,handlerfn)
}
func (s *sdkhttpServer)Start(address string)error{
	return http.ListenAndServe(address,s)
}
func (s *sdkhttpServer)ServeHTTP(w http.ResponseWriter,r *http.Request){
	ctx := NewContext(r,w)
	s.root(ctx)
}
func NewHttpServe(name string,build... filterBuilder)Server{
	handler := NewHandlerBasedOnTree()
	var root filter = handler.ServeHTTP
	for i:=len(build)-1;i>=0;i--{
		b := build[i]
		root = b(root)
	}
	return &sdkhttpServer{
		name: name,
		handler: handler,
		root: root,
	}
}