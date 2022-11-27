package queue

import (
	"context"
)

type Queue interface {
	EnQueue(ctx context.Context, val any) error
	DeQueue(ctx context.Context) (val any, err error)
}
