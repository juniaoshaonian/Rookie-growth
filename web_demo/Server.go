package main

import "net/http"

type Server interface {
	Routable
	Start(address string)error
}
type sdkHTTPServer struct {
	Name string
	root filter
	h    Handler
}
func (s *sdkHTTPServer)Route(method string,pattern string,handlefn handlerfunc)error{
	return s.h.Route(method,pattern,handlefn)
}
func (s *sdkHTTPServer)Start(address string)error{
	return http.ListenAndServe(address,s)
}
func (s *sdkHTTPServer)ServeHTTP(W http.ResponseWriter,R *http.Request){
	c := NewContext(W,R)
	s.root(c)
}
func NewsdkHttpServer(builds ...filterbuilder,name string)Server{
	handler := NewHandlerFunc()
	var root filter
	for i:=len(builds)-1;i>=0;i--{
		b := builds[i]
		root = b(root)
	}
	return &sdkHTTPServer{
		Name: name,
		root: root,
		h: handler,
	}


}