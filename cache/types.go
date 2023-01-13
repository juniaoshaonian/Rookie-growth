package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (val any, err error)
	Set(ctx context.Context, key string, val any, expiration time.Duration) error
}
