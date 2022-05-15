package main

import (
	"fmt"
	"time"
)

type filter func(c *Context)
type filterbuilder func(next filter)filter

func metricfilterbuilder(next filter)filter{
	return func(c *Context){
		start := time.Now().UnixNano()
		next(c)
		end := time.Now().UnixNano()
		fmt.Printf("耗时%d",end-start)
	}


}