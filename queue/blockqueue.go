package queue

import (
	"context"
	"sync"
)

type BlockQueue struct {
	mu         sync.Mutex // 并发安全
	maxNum     int32      //队列中的最大元素
	queue      []any      // 存放元素的容器
	num        int32
	emptyQueue chan struct{} //
	fullQueue  chan struct{}
}

//11

func (b *BlockQueue) EnQueue(ctx context.Context, val any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	// 元素入队
	b.mu.Lock()
	if b.num >= b.maxNum {
		b.mu.Unlock()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-b.fullQueue:
			b.mu.Lock()
			b.queue = append(b.queue, val)
			b.num++
			b.mu.Unlock()
			return nil
		}
	}

	select {
	case b.emptyQueue <- struct{}{}:
	default:

	}
	b.num++
	b.queue = append(b.queue, val)
	b.mu.Unlock()

}

func (b *BlockQueue) DeQueue(ctx context.Context) (val any, err error) {
	if ctx.Err() != nil {
		return val, ctx.Err()
	}
	b.mu.Lock()
	if b.num == 0 {
		b.mu.Unlock()
		select {
		case <-ctx.Done():
			return val, ctx.Err()
		case <-b.emptyQueue:
			b.mu.Lock()
			val = b.queue[0]
			b.queue = b.queue[1:]
			b.num--
			b.mu.Unlock()
			return
		}
	}
	select {
	case b.fullQueue <- struct{}{}:
	default:
	}
	val = b.queue[0]
	b.queue = b.queue[1:]
	b.num--
	b.mu.Unlock()
	return
}

type BlockQueueV2 struct {
	mu      sync.Mutex
	queue   []any
	num     int
	maxSize int
}

//
func (b *BlockQueueV2) EnQueue(ctx context.Context, val any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	b.mu.Lock()
	if b.num >= b.maxSize {
		select {
		case <-ctx.Done():

		}
	}
	b.queue = append(b.queue, val)
	b.num++
	b.mu.Unlock()
}

func (b *BlockQueueV2) DeQueue(ctx context.Context) (val any, err error) {
	if ctx.Err() != nil {
		return val, ctx.Err()
	}
	b.mu.Lock()
	val = b.queue[0]
	b.queue = b.queue[1:]
	b.num--
	b.mu.Unlock()
	return
}

type BlockQueueV1 struct {
	mu        *sync.RWMutex
	data      []any
	max       int
	num       int
	head      int
	tail      int
	fullCond  *sync.Cond
	emptyCond *sync.Cond
}

func (b *BlockQueueV1) EnQueue(ctx context.Context, val any) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	for b.IsFull() {
		ch := make(chan struct{})
		go func() {
			b.fullCond.Wait()
			select {
			case ch <- struct{}{}:
			default:
				b.fullCond.Signal()
				b.fullCond.L.Unlock()
			}
		}()
		select {
		case <-ch:
		case <-ctx.Done():
			return ctx.Err()
		}

	}
	b.mu.Lock()

	b.data[b.tail] = val
	b.tail = (b.tail + 1) % b.max
	b.num++
	b.emptyCond.Signal()
	b.mu.Unlock()
	return nil
}

func (b *BlockQueueV1) DeQueue(ctx context.Context) (val any, err error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	for b.IsEmpty() {
		ch := make(chan struct{})
		go func() {
			b.emptyCond.Wait()
			select {
			case ch <- struct{}{}:
			default:
				b.emptyCond.Signal()
				b.emptyCond.L.Unlock()
			}
		}()
		select {
		case <-ch:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

	}
	b.mu.Lock()
	val = b.data[b.head]
	b.data[b.head] = nil
	b.head = (b.head + 1) % b.max
	b.num--
	b.fullCond.Signal()
	b.mu.Unlock()
	return val, nil
}

func (b *BlockQueueV1) IsFull() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.max <= b.num {
		return true
	}
	return false
}
func (b *BlockQueueV1) IsEmpty() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.num == 0 {
		return true
	}
	return false
}

func NewBlockQueueV1(max int) *BlockQueueV1 {
	mu := sync.RWMutex{}
	return &BlockQueueV1{
		mu:        mu,
		data:      make([]any, max),
		max:       max,
		emptyCond: sync.NewCond(&mu),
		fullCond:  sync.NewCond(&mu),
	}

}
