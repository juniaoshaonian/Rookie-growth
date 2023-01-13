package rpc_service

import (
	"context"
	"io"
)

type Registry interface {
	Register(ctx context.Context, ins ServiceInstance) error
	UnRegister(ctx context.Context, ins ServiceInstance) error
	ListService(ctx context.Context, serviceName string) ([]ServiceInstance, error)
	Subscribe(serviceName string) (<-chan Event, error)
	io.Closer
}

type Event struct {
}
type ServiceInstance struct {
	Addr string
}
