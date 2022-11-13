package cache_demo

import (
	"context"
	"time"
)

type WriteThroughCache struct {
	Cache
	WriteFunc func(ctx context.Context, key string, val any) error
}

func (w *WriteThroughCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	err := w.Cache.Set(ctx, key, val, expiration)
	if err != nil {
		return err
	}
	
	return w.WriteFunc(ctx, key, val)

}
