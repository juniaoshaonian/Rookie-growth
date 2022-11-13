package cache_demo

import (
	"context"
	"log"
)

type WriteBackCache struct {
	*LocalCache
}

func NewWriteBackCache(LoadFunc func(ctx context.Context, key string, val any) error) *WriteBackCache {
	return &WriteBackCache{
		LocalCache: NewLocalCache(func(key string, val any) {
			err := LoadFunc(context.Background(), key, val)
			if err != nil {
				log.Fatal(err)
			}
		}),
	}
}

type PreLoadCache struct {
	Cache
	SentinelCache *LocalCache
}

func NewPreloadCache(c Cache, loadFunc func(ctx context.Context, key string) (any, error)) *PreLoadCache {
	return &PreLoadCache{
		Cache: c,
		SentinelCache: NewLocalCache(func(key string, val any) {
			val, err := loadFunc(context.Background(), key)
			if err == nil {
				c.Set(context.Background(), key, val, 1)
			}
		}),
	}
}
