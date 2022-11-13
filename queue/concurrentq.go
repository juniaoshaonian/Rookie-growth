package queue

import (
	"context"
	"sync"
)

type ConcurrentBlockingQueue[T any] struct {
	mutex    *sync.Mutex
	data     []T
	notFull  chan struct{}
	notEmpty chan struct{}
	maxSize  int
}

func NewConcurrentBlockingQueue[T any](maxSize int) *ConcurrentBlockingQueue[T] {
	m := &sync.Mutex{}
	return &ConcurrentBlockingQueue[T]{
		data:     make([]T, 0, maxSize),
		mutex:    m,
		notFull:  make(chan struct{}, 1),
		notEmpty: make(chan struct{}, 1),
		maxSize:  maxSize,
	}
}

func (c *ConcurrentBlockingQueue[T]) EnQueue(ctx context.Context, data T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	c.mutex.Lock()
	for c.IsFull() {
		c.mutex.Unlock()
		select {
		case <-c.notFull:
			c.mutex.Lock()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	c.data = append(c.data, data)
	if len(c.data) == 1 {
		c.notEmpty <- struct{}{}
	}
	c.mutex.Unlock()
	// 没有人等 notEmpty 的信号，这一句就会阻塞住
	return nil
}

func (c *ConcurrentBlockingQueue[T]) DeQueue(ctx context.Context) (T, error) {
	if ctx.Err() != nil {
		var t T
		return t, ctx.Err()
	}
	c.mutex.Lock()
	for c.IsEmpty() {
		c.mutex.Unlock()
		select {
		case <-c.notEmpty:
			c.mutex.Lock()
		case <-ctx.Done():
			var t T
			return t, ctx.Err()
		}
	}
	t := c.data[0]
	c.data = c.data[1:]
	if len(c.data) == c.maxSize-1 {
		c.notFull <- struct{}{}
	}
	c.mutex.Unlock()
	// 没有人等 notFull 的信号，这一句就会阻塞住
	return t, nil
}

func (c *ConcurrentBlockingQueue[T]) IsFull() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return len(c.data) == c.maxSize
}

func (c *ConcurrentBlockingQueue[T]) IsEmpty() bool {
	return len(c.data) == 0
}

func (c *ConcurrentBlockingQueue[T]) Len() uint64 {
	return uint64(len(c.data))
}
