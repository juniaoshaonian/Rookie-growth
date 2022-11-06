package cache_demo

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"time"
)

type RedisCache struct {
	client    redis.Cmdable
	onEvicted func(key string, val any)
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	return r.client.Get(ctx, key).Result()

}

func (r *RedisCache) Set(ctx context.Context, key string, val any, expairtion time.Duration) error {
	res, err := r.client.Set(ctx, key, val, expairtion).Result()
	if err != nil {
		return err
	}
	if res != "OK" {
		return errors.New("cache 设置键值对失败")
	}
	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}

func (r *RedisCache) OnEvicted(f func(key string, val any)) {
	r.onEvicted = f
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{
		client: client,
	}
}
