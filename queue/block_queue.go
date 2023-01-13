package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

type BlockQueueV3 struct {
	mu        *sync.RWMutex
	data      []any
	max       int
	num       int
	emptyCond *Cond
	fullCond  *Cond
}

func (b *BlockQueueV3) EnQueue(ctx context.Context, val any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	b.mu.Lock()
	for b.IsFull() {
		err := b.fullCond.WaitTimeOut(ctx)
		if err != nil {
			return err
		}
	}
	b.mu.Lock()
	b.data = append(b.data, val)
	b.num++
	b.emptyCond.Broadcast()
	b.mu.Unlock()
	return nil
}

func (b *BlockQueueV3) DeQueue(ctx context.Context) (val any, err error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	b.mu.Lock()
	for b.IsEmpty() {
		b.emptyCond.WaitTimeOut(ctx)
		if err != nil {
			return nil, err
		}
	}
	b.mu.Lock()
	val = b.data[0]
	b.data = b.data[1:]
	b.num--
	b.fullCond.Broadcast()
	b.mu.Unlock()
	return
}

func (b *BlockQueueV3) IsFull() bool {
	if b.num >= b.max {
		return true
	}
	return false
}
func (b *BlockQueueV3) IsEmpty() bool {
	if b.num <= 0 {
		return true
	}
	return false
}
func NewBlockQueueV3(max int) *BlockQueueV3 {
	l := &sync.RWMutex{}
	return &BlockQueueV3{
		data:      make([]any, 0, max),
		mu:        l,
		max:       max,
		emptyCond: NewCond(l),
		fullCond:  NewCond(l),
	}
}

type Cond struct {
	L sync.Locker
	n unsafe.Pointer
}

func NewCond(l sync.Locker) *Cond {
	c := &Cond{
		L: l,
	}
	ch := make(chan struct{})
	c.n = unsafe.Pointer(&ch)
	return c
}

func (c *Cond) WaitTimeOut(ctx context.Context) error {
	n := c.NotifyChan()
	c.L.Unlock()
	select {
	case <-n:
		c.L.Lock()
	case <-ctx.Done():
		c.L.Lock()
		return ctx.Err()
	}

	return nil
}

func (c *Cond) NotifyChan() <-chan struct{} {
	ptr := atomic.LoadPointer(&c.n)
	return *((*chan struct{})(ptr))
}

func (c *Cond) Broadcast() {
	// 加载出来channel
	ch := make(chan struct{})
	oldPtr := atomic.SwapPointer(&c.n, unsafe.Pointer(&ch))
	close(*((*chan struct{})(oldPtr)))
}
