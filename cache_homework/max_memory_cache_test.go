package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 测试正常情况
func TestMaxMemoryCache_Set(t *testing.T) {
	// 模拟没有淘汰的情况
	mcache := NewMaxMemoryCache(100)
	testcases := []node{
		{
			"key1",
			[]byte("value1"),
		},
		{
			"key2",
			[]byte("value2"),
		},
		{
			"key3",
			[]byte("value3"),
		},
	}
	for _, tc := range testcases {
		mcache.Set(context.Background(), tc.key, tc.val, 1)
	}
	assert.Equal(t, 3, len(mcache.datas))
	for _, tc := range testcases {
		val, err := mcache.Get(context.Background(), tc.key)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tc.val, val)
	}

	// 测试队列满然后淘汰的情况,key4太大会导致cache淘汰key1，和key2
	testcases = append(testcases, node{
		"key4",
		[]byte("value4value4value4value4value4value4value4"),
	})
	mcache.Set(context.Background(), "key4", []byte("value4value4value4value4value4value4value4"), 1)
	assert.Equal(t, 2, len(mcache.datas))
	for _, tc := range testcases {
		val, err := mcache.Get(context.Background(), tc.key)
		if tc.key == "key1" || tc.key == "key2" {
			assert.Equal(t, errors.New(fmt.Sprintf("cache can not found key:%s", tc.key)), err)
			continue
		}
		assert.Equal(t, tc.val, val)
	}
	// 测试修改队列已存在的值使缓存满淘汰最老的键,将key4的值改大会将剩下的key3也给淘汰出去
	testcases[3].val = []byte("value4value4value4value4value4value4value4value4value4value4value4value4")
	mcache.Set(context.Background(), "key4", []byte("value4value4value4value4value4value4value4value4value4value4value4value4"), 1)
	for _, tc := range testcases {
		val, err := mcache.Get(context.Background(), tc.key)
		if tc.key == "key1" || tc.key == "key2" || tc.key == "key3" {
			assert.Equal(t, errors.New(fmt.Sprintf("cache can not found key:%s", tc.key)), err)
			continue
		}
		assert.Equal(t, tc.val, val)
	}
}

type node struct {
	key string
	val []byte
}
