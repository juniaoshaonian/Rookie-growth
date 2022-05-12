package main
type matchFunc func()
type node struct {
	children  []*node
	handler handlerfunc
	ma
}