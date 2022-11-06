package cache_demo

import (
	"cache_demo/mocks"
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRedisCache_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testcases := []struct {
		name       string
		mock       func() redis.Cmdable
		key        string
		val        any
		expiration time.Duration
		wantErr    error
	}{
		{
			name: "Return OK",
			mock: func() redis.Cmdable {
				res := mocks.NewMockCmdable(ctrl)
				cmd := redis.NewStatusCmd(nil)
				cmd.SetVal("OK")
				res.EXPECT().Set(gomock.Any(), "key1", "value1", time.Minute).Return(cmd)
				return res
			},
			key:        "key1",
			val:        "value1",
			expiration: time.Minute,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cmdable := tc.mock()
			client := NewRedisCache(cmdable)
			err := client.Set(context.Background(), tc.key, tc.val, tc.expiration)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
