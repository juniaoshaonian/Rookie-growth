package cache_demo

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any, expairtion time.Duration) error
	Delete(ctx context.Context, key string) error
	OnEvicted(func(key string, val any))
}
