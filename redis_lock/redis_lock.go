package redis_lock

import (
	"context"
	_ "embed"
	"errors"
	"github.com/go-redis/redis/v9"
	"github.com/google/uuid"
	"redis_lock/retry"
	"time"
)

type Client struct {
	client redis.Cmdable
}

var (
	//go:embed lua/unlock.lua
	luaUnLock string
	//go:embed lua/refresh.lua
	luaRefresh string
	//go:embed lua/lock.lua
	luaLock                string
	ErrNoLockHeld          error = errors.New("未持有锁")
	ErrFailedToPreemptLock       = errors.New("rlock: 抢锁失败")
)

func (c *Client) TryLock(ctx context.Context, key string, expiration time.Duration, timeout time.Duration) (*Lock, error) {
	val := uuid.New().String()
	timeoutctx, cancel := context.WithTimeout(ctx, timeout)
	ok, err := c.client.SetNX(timeoutctx, key, val, expiration).Result()
	cancel()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrNoLockHeld
	}
	return NewLock(key, val, expiration, c), nil
}

type Lock struct {
	key        string
	val        string
	expiration time.Duration
	client     *Client
	unlock     chan struct{}
}

func NewLock(key string, val string, expiration time.Duration, c *Client) *Lock {
	return &Lock{
		key:        key,
		val:        val,
		expiration: expiration,
		client:     c,
	}
}

func (l *Lock) UnLock(ctx context.Context) error {
	res, err := l.client.client.Eval(ctx, luaUnLock, []string{l.key}, l.val).Result()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrNoLockHeld
	}
	return nil
}

func (l *Lock) Refresh(ctx context.Context) error {
	res, err := l.client.client.Eval(ctx, luaRefresh, []string{l.key}, l.val, l.expiration.Seconds()).Result()
	if err != nil {
		return err
	}
	if res != 1 {
		return ErrNoLockHeld
	}
	return nil
}

func (l *Lock) AutoRefresh(ctx context.Context, interval time.Duration, timeout time.Duration) error {
	//续约机制
	//1. 每隔一段时间续约一次
	//2. 遇到超时错误立即重试
	//3. 解锁时需要退出循环
	ticker := time.NewTicker(interval)
	retry := make(chan struct{}, 1)
	defer func() {
		ticker.Stop()
		close(retry)
	}()

	for {
		select {
		case <-ticker.C:
			timeoutctx, cancel := context.WithTimeout(ctx, timeout)
			err := l.Refresh(timeoutctx)
			cancel()
			if err != nil && err != context.DeadlineExceeded {
				return err
			}
			if err == context.DeadlineExceeded {
				retry <- struct{}{}
				continue
			}
		case <-retry:
			timeoutctx, cancel := context.WithTimeout(ctx, timeout)
			err := l.Refresh(timeoutctx)
			cancel()
			if err != nil && err != context.DeadlineExceeded {
				return err
			}
			if err == context.DeadlineExceeded {
				retry <- struct{}{}
				continue
			}
		case <-l.unlock:
			return nil
		}
	}
}

// 如果需要重试，我们
func (c *Client) Lock(ctx context.Context, key string, expiration time.Duration, retry retry.Retry, timeout time.Duration) (*Lock, error) {
	val := uuid.New().String()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			timeoutctx, cancel := context.WithTimeout(ctx, timeout)
			res, err := c.client.Eval(timeoutctx, luaLock, []string{key}, val, expiration.Seconds()).Result()
			cancel()
			if err != nil && err != context.DeadlineExceeded {
				return nil, err
			}
			if res == "OK" {
				return &Lock{
					key:        key,
					val:        val,
					expiration: expiration,
				}, nil
			}
			if res == 0 {
				return &Lock{
					key:        key,
					val:        val,
					expiration: expiration,
				}, nil
			} else {
				return nil, ErrFailedToPreemptLock
			}
			ok, interval := retry.Next()
			if !ok {
				return nil, errors.New("重试次数过多")
			}
			time.Sleep(interval)
		}
	}
}
