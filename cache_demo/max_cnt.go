package cache_demo

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

type MaxCountCache struct {
	MaxCnt int32
	cnt    int32
	Cache
}

func NewMaxCountCache(MaxCnt int32) *MaxCountCache {
	m := &MaxCountCache{MaxCnt: MaxCnt}
	f := func(key string, val any) {
		atomic.AddInt32(&m.cnt, -1)
	}
	m.OnEvicted(f)
	return m
}

func (m *MaxCountCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	kk := atomic.AddInt32(&m.cnt, 1)
	if kk > m.MaxCnt {
		atomic.AddInt32(&m.cnt, -1)
		return errors.New("cache is full")
	}
	return m.Cache.Set(ctx, key, val, expiration)
}
