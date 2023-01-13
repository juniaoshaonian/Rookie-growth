package rpc_demo2

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/resolver"
	"rpc"
	"rpc_demo2/mocks"
	"rpc_demo2/register"
	"testing"
	"time"
)

func TestNewResolverBuilder_ResolverNow(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testcases := []struct {
		name    string
		r       func() register.Registry
		wantVal resolver.State
	}{
		{
			name: "",
			r: func() register.Registry {
				r := mocks.NewMockRegistry(ctrl)
				r.EXPECT().ListService(gomock.Any(), gomock.Any()).Return([]register.ServiceInstance{
					{
						Addr: "service-1",
					},
					{
						Addr: "service-2",
					},
				}, nil)
				return r
			},
			wantVal: resolver.State{
				Addresses: []resolver.Address{
					resolver.Address{
						Addr: "service-1",
					},
					resolver.Address{
						Addr: "service-2",
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			cc := &MockConnClient{}
			g := rpc.grpcResolver{
				r:      tc.r(),
				target: resolver.Target{},
				cc:     cc,
			}
			g.resolve()
			assert.Equal(t, tc.wantVal, cc.state)

		})
	}

}

type MockConnClient struct {
	resolver.ClientConn
	state resolver.State
}

func (m *MockConnClient) UpdateState(s resolver.State) error {
	m.state = s
	return nil
}

func TestGrpcResolver_watch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// 是否能退出，是否可以正常接收信号
	testcases := []struct {
		name      string
		mock      func() (register.Registry, chan register.Event)
		wantState resolver.State
		wantErr   error
	}{
		{
			name: "",
			mock: func() (register.Registry, chan register.Event) {
				r := mocks.NewMockRegistry(ctrl)
				ch := make(chan register.Event)
				r.EXPECT().SubScribe(gomock.Any()).Return(ch, nil)
				r.EXPECT().ListService(gomock.Any(), gomock.Any()).
					Return([]register.ServiceInstance{
						{
							Addr: "test-1",
						},
					}, nil)
				return r, ch
			},
			wantState: resolver.State{
				Addresses: []resolver.Address{
					{
						Addr: "test-1",
					},
				},
			},
		},
	}
	for _, tc := range testcases {
		cc := &MockConnClient{}
		r, ch := tc.mock()
		closech := make(chan struct{})
		rs := &rpc.grpcResolver{
			r:       r,
			cc:      cc,
			closech: closech,
		}
		err := rs.watch()
		assert.Equal(t, tc.wantErr, err)
		ch <- register.Event{}
		time.Sleep(time.Second)
		assert.Equal(t, tc.wantState, cc.state)
		rs.Close()
		_, ok := <-closech
		assert.False(t, ok)

	}

}
