package register

import "context"

type Registry interface {
	Register(ctx context.Context, ins ServiceInstance) error
	UnRegiter(ctx context.Context, ins ServiceInstance) error
	ListService(ctx context.Context, serviceName string) ([]*ServiceInstance, error)
	SubScribe(serviceName string) (<-chan Event, error)
	Close() error
}

type ServiceInstance struct {
	ServiceName string
	Addr        string
}

type Event struct {
	Type EventType
}
type EventType int

const (
	EventTypeUnknown EventType = iota
	EventTypeAdd
	EventTypeDelete
	EventTypeUpdate
	// EventTypeErr
)
