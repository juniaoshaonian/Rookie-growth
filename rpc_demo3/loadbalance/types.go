package loadbalance

import (
	"context"
	"google.golang.org/grpc/resolver"
)

type Filter func(ctx context.Context, addr resolver.Address) bool

func GroupFilter(ctx context.Context, addr resolver.Address) bool {
	group := ctx.Value("group")
	if group == nil {
		return true
	}
	if group != addr.Attributes.Value("group") {
		return false
	}
	return true
}
