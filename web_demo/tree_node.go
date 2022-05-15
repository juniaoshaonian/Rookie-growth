package main

import "strings"

const (
	NodeRoot = iota
	NodeAny
	NodeParam
	NodeStatic
)
const any = "*"
type matchFunc func(string,*Context)bool
type Node struct {
	pattern string
	node_type int
	children []*Node
	handler handlerfunc
	m matchFunc
}
func NewRootNode(method string)*Node {
	return &Node{
		pattern: method,
		node_type: NodeRoot,
		children: make([]*Node,0,4),
		m : func(path string,c *Context)bool{
			panic("nevery call me")
		},
	}
}
func NewStaicNode(pattern string)*Node {
	return &Node{
		pattern: pattern,
		node_type: NodeStatic,
		children: make([]*Node,0,4),
		m: func(path string,c *Context)bool{
			if path == pattern && pattern != "*" {
				return true
			}
			return false
		},
	}
}
func NewAnyNode()*Node{
	return &Node{
		pattern: any,
		node_type: NodeAny,
		//children: make([]*Node,0,4),
		m: func(path string,c *Context)bool{
			return true
		},
	}
}
func NewParamNode(pattern string)*Node{
	p := pattern[1:]
	return &Node{
		pattern: pattern,
		node_type: NodeParam,
		m: func(path string,c *Context)bool{
			if c != nil {
				c.param[p] = path
			}
			return p != any
		},
		children: make([]*Node,0,4),
	}
}
func NewNode(pattern string)*Node{
	if pattern == "*" {
		return NewAnyNode()
	}
	if strings.HasPrefix(pattern,":") {
		return NewParamNode(pattern)
	}
	return NewStaicNode(pattern)
}
func (n *Node)CreatesubTree(paths []string,handlefn handlerfunc){
	this := n
	for _,path := range paths {
		new_node := NewNode(path)
		this.children = append(this.children,new_node)
		this = new_node
	}
	this.handler = handlefn
}
