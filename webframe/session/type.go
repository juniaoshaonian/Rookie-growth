package session

import (
	"context"
	"net/http"
)

type Session interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, val string) error
	ID() string
}

// 决定怎样存储
type Store interface {
	Get(ctx context.Context, id string) (Session, error)
	Genrate(ctx context.Context, id string) (Session, error)
	Remove(ctx context.Context, id string) error
	Reflash(ctx context.Context, id string) error
}

type Propagator interface {
	// 将session写到response里面
	Inject(id string, resp http.ResponseWriter) error
	Extract(req *http.Request) (string, error)
	Remove(resp http.ResponseWriter) error
}
