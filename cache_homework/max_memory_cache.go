package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type MaxMemoryCache struct {
	Cache
	max       int64
	used      int64
	head      *DataNode
	tail      *DataNode
	datas     map[string]*DataNode
	onEvicted func(key string, val []byte)
}

func (l *MaxMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	data, ok := l.datas[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("cache can not found key:%s", key))
	}
	l.DeleteNode(data)
	l.PutBehindHead(data)
	return data.val, nil
}

func (l *MaxMemoryCache) PutBehindHead(node *DataNode) {
	next := l.head.next
	l.head.next = node
	node.pre = l.head
	node.next = next
	next.pre = node
}
func (l *MaxMemoryCache) DeleteNode(node *DataNode) {
	node.pre.next = node.next
	node.next.pre = node.pre
}
func (l *MaxMemoryCache) Set(ctx context.Context, key string, val []byte, expiration time.Duration) error {
	data, ok := l.datas[key]
	if ok {
		l.Delete(context.Background(), key)
		data.val = val
		l.PutBehindHead(data)
		l.datas[key] = data
		l.used = l.used + int64(len(key)) + int64(len(val)) + int64(16)
	} else {
		new_node := &DataNode{
			key: key,
			val: val,
		}
		l.datas[key] = new_node
		l.PutBehindHead(new_node)
		l.used = l.used + int64(len(key)) + int64(len(val)) + int64(16)
	}

	for l.used > l.max {
		lastNode := l.tail.pre
		l.Delete(ctx, lastNode.key)
	}
	return nil
}
func (l *MaxMemoryCache) Delete(ctx context.Context, key string) error {
	val, ok := l.datas[key]
	if !ok {
		return errors.New("KEY NOT FOUND")
	}
	l.used = l.used - int64(len(key)) - int64(len(val.val)) - 2*int64(8)
	delete(l.datas, key)
	l.DeleteNode(val)
	if l.onEvicted != nil {
		l.onEvicted(key, val.val)
	}
	return nil
}

type DataNode struct {
	pre  *DataNode
	next *DataNode
	key  string
	val  []byte
}

func NewMaxMemoryCache(max int64) *MaxMemoryCache {
	m := &MaxMemoryCache{
		max:   max,
		datas: make(map[string]*DataNode, max),
	}
	m.head = &DataNode{}
	m.tail = &DataNode{}
	m.head.next = m.tail
	m.tail.pre = m.head
	return m
}
