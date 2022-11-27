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
