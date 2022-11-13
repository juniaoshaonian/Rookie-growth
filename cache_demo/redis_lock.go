package cache_demo

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"time"
)

var (
	ErrFailedToPreemptLock = errors.New("rlock: 抢锁失败")
	// ErrLockNotHold 一般是出现在你预期你本来持有锁，结果却没有持有锁的地方
	// 比如说当你尝试释放锁的时候，可能得到这个错误
	// 这一般意味着有人绕开了 rlock 的控制，直接操作了 Redis
	ErrLockNotHold = errors.New("rlock: 未持有锁")
)

type Client struct {
	client redis.Cmdable
}

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration) error {
	value := uuid.New()
	c.client.Eval()
	ok, err := c.client.SetNX(ctx, key, value, expiration).Result()
	if err != nil {
		return err
	}
	if !ok {
		return ErrFailedToPreemptLock
	}
	return nil
}

func (c *Client) Unlock(ctx context.Context, key string) error {
	res, err := c.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrLockNotHold
	}
}
