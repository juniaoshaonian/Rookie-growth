package cache_demo

import (
	"cache_demo/errs"
	"context"
	"log"
)

type ReadThroughCache struct {
	Cache
	LoadFunc func(ctx context.Context, key string) (val any, err error)
}

func (r *ReadThroughCache) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err != nil && err != errs.ErrKeyNotFound {
		return nil, err
	}
	if err == errs.ErrKeyNotFound {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			err = r.Set(ctx, key, val, 1)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return val, err
}
