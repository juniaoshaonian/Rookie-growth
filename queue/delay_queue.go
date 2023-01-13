package queue

import (
	"context"
	"sync"
	"time"
)

type DelayQueue[T DelayAble] struct {
	mu        *sync.RWMutex
	p         PriorityQueue[T]
	fullCond  *Cond
	emptyCond *Cond
}

func (d *DelayQueue[T]) EnQueue(ctx context.Context, val T) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		d.mu.Lock()
		err := d.p.Enqueue(val)
		switch err {
		case nil:

			d.fullCond.Broadcast()
			d.mu.Unlock()
			return nil
		case ErrOutOfCapacity:
			ch := d.emptyCond.NotifyChan()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ch:

			}
		default:
			d.mu.Unlock()
			return err
		}
	}
}

func (d *DelayQueue[T]) DeQueue(ctx context.Context) (DelayAble, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		d.mu.Lock()
		head, err := d.p.Peek()
		d.mu.Unlock()
		switch err {
		case nil:
			if head.Delay() <= 0 {
				d.mu.Lock()
				val, err := d.p.Dequeue()
				d.mu.Unlock()
				if err != nil {
					d.fullCond.Broadcast()
				}
				return val, err
			}
			timer := time.NewTimer(head.Delay())
			ch := d.emptyCond.NotifyChan()
			select {
			case <-timer.C:
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-ch:
			}
		case ErrEmptyQueue:
			ch := d.emptyCond.NotifyChan()
			select {
			case <-ctx.Done():
				return nil, err
			case <-ch:
			}
		default:
			return nil, err

		}

	}
}

type DelayAble interface {
	Delay() time.Duration
}
