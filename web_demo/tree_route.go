package main

import (
	"errors"
	"net/http"
	"sort"
	"strings"
)
var ErrorInvalidMethod error = errors.New("err method")
var ErrorInvalidPattern error = errors.New("error pattern")
var supportMethod []string =[]string{
	http.MethodGet,
	http.MethodPost,
	http.MethodDelete,
	http.MethodPut,
}
type HandlerBasedOnTree struct {
	Forest map[string]*Node
}

func NewHandlerFunc() Handler{
	xx := &HandlerBasedOnTree{}
	for _,method:=range supportMethod{
		xx.Forest[method] =  NewRootNode(method)
	}
	return xx
}
func (h *HandlerBasedOnTree)Route(method string,pattern string,handlefn handlerfunc)error {
	root,ok := h.Forest[method]
	if !ok {
		return ErrorInvalidMethod
	}
	pattern = strings.Trim(pattern,"/")
	err := h.InvalidError(pattern)
	if err != nil {
		return err
	}
	paths := strings.Split(pattern,"/")
	for index,path :=range paths {
		matchchild,found := h.findchild(root,path,nil)
		if found {
			root = matchchild
		}else  {
			root.CreatesubTree(paths[index:],handlefn)
			return nil
		}
	}
	root.handler = handlefn
	return nil
}
func (h *HandlerBasedOnTree)findchild(cur *Node,path string,c *Context)(*Node,bool){
	candidates := make([]*Node,0,2)
	for _,child := range cur.children {
		if child.m(path,c) {
			candidates = append(candidates,child)
		}
	}
	if len(candidates)==0 {
		return nil,false
	}
	sort.Slice(candidates,func(i,j int)bool{
		return candidates[i].node_type < candidates[j].node_type
	})
	return candidates[len(candidates)-1],true
}
func (h *HandlerBasedOnTree)InvalidError(path string)error{
	pos := strings.Index(path,"*")
	if pos != -1 && pos != len(path)-1 {
		return ErrorInvalidPattern
	}
	if pos != -1 && path[pos-1] != '/' {
		return ErrorInvalidPattern
	}
	return nil
}
func (h *HandlerBasedOnTree)ServeHTTP(c *Context){
	hand,ok:=h.findroute(c,c.R.URL.Path)
	if !ok {
		c.W.WriteHeader(http.StatusNotFound)
		return
	}
	hand(c)
}
func (h *HandlerBasedOnTree)findroute(c *Context,path string)(handlerfunc,bool){
	root,ok := h.Forest[c.R.Method]
	if !ok {
		return nil,false
	}
	path = strings.Trim(path,"/")
	paths := strings.Split(path,"/")
	for _,p := range paths{
		matchchild,found := h.findchild(root,p,c)
		if found {
			root = matchchild
		}else {
			return nil,false
		}
	}
	if root.handler == nil {
		return nil,false
	}
	return root.handler,true
}