package queue

import (
	"context"
	"sync/atomic"
	"unsafe"
)

type LinkQueue[T any] struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func (l *LinkQueue[T]) EnQueue(ctx context.Context, val T) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	new_node := &DataVal[T]{
		val: val,
	}

	newNodePtr := unsafe.Pointer(new_node)
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		oldPtr := atomic.LoadPointer(&l.tail)
		if atomic.CompareAndSwapPointer(&l.tail, oldPtr, newNodePtr) {
			tail := (*DataVal[T])(oldPtr)
			atomic.StorePointer(&tail.next, oldPtr)
			return nil
		}

	}
}

func (l *LinkQueue[T]) DeQueue(ctx context.Context) (T, error) {
	var t T
	if ctx.Err() != nil {

		return t, ctx.Err()
	}
	for {
		if ctx.Err() != nil {
			return t, ctx.Err()
		}
		headPtr := atomic.LoadPointer(&l.head)
		head := (*DataVal[T])(headPtr)
		tailPtr := atomic.LoadPointer(&l.tail)
		tail := (*DataVal[T])(tailPtr)

		if head == tail {
			continue
		}
		next := atomic.LoadPointer(&head.next)
		if atomic.CompareAndSwapPointer(&l.head, headPtr, next) {
			nextPtr := (*DataVal[T])(next)
			return nextPtr.val, nil
		}
	}

}

type DataVal[T any] struct {
	next unsafe.Pointer
	val  T
}
