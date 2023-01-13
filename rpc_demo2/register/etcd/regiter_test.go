package etcd

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"rpc_demo2/register"
	"rpc_demo2/register/etcd/mocks"
	"testing"
)

func TestRegiter_Subscribe(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testcases := []struct {
		name      string
		mock      func() (clientv3.Watcher, chan clientv3.WatchResponse)
		Response  clientv3.WatchResponse
		wantEvent register.Event
		wantErr   error
	}{
		{
			name: "subscribe",
			mock: func() (clientv3.Watcher, chan clientv3.WatchResponse) {
				watcher := mocks.NewMockWatcher(ctrl)
				ch := make(chan clientv3.WatchResponse)
				watcher.EXPECT().Watch(gomock.Any(), gomock.Any(), gomock.Any()).Return(ch)
				return watcher, ch

			},
			Response: clientv3.WatchResponse{
				Events: []*clientv3.Event{
					{
						Type: clientv3.EventTypePut,
						Kv:  api.,
					},
				},
			},
			wantEvent: register.Event{
				Type: register.EventTypeAdd,
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			// 测试是否可以正常接收数据
			watcher, ch := tc.mock()
			r := EtcdRegister{
				client: &clientv3.Client{
					Watcher: watcher,
				},
			}
			event, err := r.SubScribe("xxxx")
			assert.Equal(t, tc.wantErr, err)
			stopch := make(chan struct{})
			go func() {
				select {
				case e := <-event:
					assert.Equal(t, tc.wantEvent, e)
					stopch <- struct{}{}

				}
			}()
			ch <- tc.Response
			<-stopch

		})
	}
}
