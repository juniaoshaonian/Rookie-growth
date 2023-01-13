package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type LocalCache struct {
	data   map[string]*Item
	mu     *sync.RWMutex
	stopCh chan struct{}
}

func (l *LocalCache) Get(ctx context.Context, key string) (val any, err error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	l.mu.RLock()
	value, ok := l.data[key]
	l.mu.RUnlock()
	if !ok {
		return nil, errors.New("key not found")
	}
	if value.deadline.Before(time.Now()) {
		l.mu.Lock()
		value, ok = l.data[key]
		if ok {
			delete(l.data, key)
		}
		l.mu.Unlock()
		return nil, errors.New("key not found")
	}
	return value, nil

}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.mu.Lock()
	l.data[key] = NewItem(val, expiration)
	l.mu.Unlock()
	return nil

}

func NewLocalCache(interval time.Duration) *LocalCache {
	localCache := &LocalCache{
		data: make(map[string]*Item, 16),
	}
	stopChan := make(chan struct{}, 1)
	localCache.stopCh = stopChan
	go func() {
		tick := time.NewTicker(interval)
		defer tick.Stop()
		for {
			select {
			case <-stopChan:
				return
			case <-tick.C:
				localCache.mu.Lock()
				for key, val := range localCache.data {
					if val.deadline.Before(time.Now()) {
						delete(localCache.data, key)
					}
				}
				localCache.mu.Unlock()
			}
		}

	}()
	return localCache
}

type Item struct {
	val      any
	deadline time.Time
}

func NewItem(val any, expiration time.Duration) *Item {
	return &Item{
		val:      val,
		deadline: time.Now().Add(expiration),
	}
}
