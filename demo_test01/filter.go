package main

import (
	"fmt"
	"time"
)

type filterBuilder func(h filter)filter
type filter func(ctx *Context)

func MetricFilter(n filter)filter {
	return func(ctx *Context) {
		start := time.Now().UnixNano()
		n(ctx)
		end := time.Now().UnixNano()
		fmt.Printf("共耗时%d",end-start)
	}
}