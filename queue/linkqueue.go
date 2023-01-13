package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"unsafe"
)

type ConcurrentLinkedQueue[T any] struct {
	head  unsafe.Pointer
	tail  unsafe.Pointer
	mutex sync.RWMutex
}

func (c *ConcurrentLinkedQueue[T]) EnQueue(ctx context.Context, data T) error {

	newNode := &node[T]{
		val: data,
	}

	tail := atomic.LoadPointer(&c.tail)
	atomic.CompareAndSwapPointer(&c.tail,
		tail,
		unsafe.Pointer(newNode))

	// 修改 tail.next
	tail.next = atomic.CompareAndSwapPointer(&)
	return nil

}

func (c *ConcurrentLinkedQueue[T]) DeQueue(ctx context.Context) (T, error) {

	for {
		head := atomic.LoadPointer(&c.head)

		if atomic.CompareAndSwapPointer(&c.head,head,head.) {
			return 
		}



	}



	return head.val.(T), nil
}

func (c *ConcurrentLinkedQueue[T]) IsFull() bool {
	// TODO implement me
	panic("implement me")
}

func (c *ConcurrentLinkedQueue[T]) IsEmpty() bool {
	// TODO implement me
	panic("implement me")
}

func (c *ConcurrentLinkedQueue[T]) Len() uint64 {
	// TODO implement me
	panic("implement me")
}

type node struct {
	next *node
	val  any
}
