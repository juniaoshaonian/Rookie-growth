package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server interface {
	Routable
	Start(address string)error
	Shutdown(ctx context.Context)error

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
func NewsdkHttpServer(name string,builds... filterbuilder,)Server{
	handler := NewHandlerFunc()
	var root filter = handler.ServeHTTP
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
func (s *sdkHTTPServer)Shutdown(ctx context.Context)error{
	fmt.Printf("#{s.Name} shutdown...\n")
	time.Sleep(time.Second)
	fmt.Printf("#{s.Name} shutdown!!!\n")
	return nil

}