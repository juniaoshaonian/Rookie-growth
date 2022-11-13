package queue

import "context"

type Queue[T any] interface {
	Enqueue(ctx context.Context, val T) error
	Dequeue(ctx context.Context) (T, error)
}
