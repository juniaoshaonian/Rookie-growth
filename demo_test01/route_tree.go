package main

import (
	"net/http"
	"strings"
)

type HttpHandlerbasedonTree struct {
	root *node
}
type node struct {
	 pattern string
	 children []*node
	 handlefn handlerfunc
}
func (h *HttpHandlerbasedonTree)Route(method string,Pattern string,handlerfn handlerfunc){
	Pattern = strings.Trim(Pattern,"/")
	paths := strings.Split(Pattern,"/")
	cur := h.root
	for index,path := range paths{
		if matchild,found := cur.findchild(path);found{
			cur = matchild
		}else {
			cur.CreatesubTree(cur,paths[index:],handlerfn)
			return
		}
	}
}
func (h *HttpHandlerbasedonTree)ServeHTTP(ctx *Context){
	xx,found := h.FindFunc(ctx.R.URL.Path)
	if !found {
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("not found"))
	}
	xx(ctx)
}
func (h *HttpHandlerbasedonTree)FindFunc(Pattern string)(handlerfunc,bool){
	Pattern = strings.Trim(Pattern,"/")
	paths := strings.Split(Pattern,"/")
	cur := h.root
	for _,path := range paths {
		if matchhild,found := cur.findchild(path);found {
			cur =  matchhild
		}else {
			return nil,false
		}
	}
	if cur.handlefn == nil {
		return nil,false
	}
	return cur.handlefn,true
}
func (n *node)findchild(path string)(*node,bool){
	cur := n
	for _,child := range cur.children {
		if child.pattern == path {
			return child,true
		}
	}
	return nil,false
}
func (n *node)CreatesubTree(cur *node,paths []string,handlerfn handlerfunc){
	this := cur
	for _,path := range paths {
		new_node := Newnode(path)
		cur.children = append(cur.children,new_node)
		this = new_node
	}
	this.handlefn = handlerfn
}
func Newnode(pattern string)*node{
	return &node{
		pattern: pattern,
		children: make([]*node,0),
	}
}
func NewHandlerBasedOnTree() *HttpHandlerbasedonTree{
	return &HttpHandlerbasedonTree{}
}