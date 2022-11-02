package cache_demo

import (
	"context"
	"sync"
	"time"
)

type LocalCache struct {
	close     chan struct{}
	m         sync.Map
	closeOnce sync.Once
}

func (l *LocalCache) Get(ctx context.Context, key string) (any, error) {
	val, ok := l.m.Load(key)
	if !ok {
		
	}
}

func (l *LocalCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	l.m.Store(key, &item{
		val:      val,
		deadline: time.Now().Add(expiration),
	})
}

func (l *LocalCache) Delete(ctx context.Context, key string) error {
	//TODO implement me
	panic("implement me")
}

func NewLocalCache() *LocalCache {
	res := &LocalCache{
		close: make(chan struct{}, 1),
	}
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				res.m.Range(func(key, value any) bool {
					itm := value.(*item)
					if itm.deadline.Before(time.Now()) {
						res.m.Delete(key)
					}

					return true
				})
			case <-res.close:

			}
		}
	}()
}

func (l *LocalCache) Close() error {
	l.closeOnce.Do(func() {
		l.close <- struct{}{}
		close(l.close)
	})

	return nil
}

type item struct {
	val      any
	deadline time.Time
}
