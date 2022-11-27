package rpc_demo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitClientProxy(t *testing.T) {
	testcases := []struct {
		name     string
		proxy    *MockProxy
		wantReq  *Request
		wantResp *GetUserResp
		wantErr  error
	}{
		{

			name: "initproxy",
			proxy: &MockProxy{
				data: []byte(`{"name":"tom"}`),
			},
			wantResp: &GetUserResp{
				Name: "tom",
			},
			wantReq: &Request{
				ServiceName: "user_service",
				MethodName:  "GetUserId",
				Arg:         []byte(`{"Id":13}`),
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			us := &UserServiceClient{}
			err := InitClientProxy(us, tc.proxy)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := us.GetUserId(context.Background(), &GetUserReq{Id: 13})
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantReq, tc.proxy.req)
			assert.Equal(t, tc.wantResp, resp)
		})
	}

}

// mockproxy
type MockProxy struct {
	// 在这写数据的原因，是为了可以定制data
	data []byte

	req *Request
}

func (m *MockProxy) Invoke(ctx context.Context, req *Request) (*Response, error) {
	m.req = req
	return &Response{
		Data: m.data,
	}, nil
}

type UserServiceClient struct {
	GetUserId func(ctx context.Context, req *GetUserReq) (resp *GetUserResp, err error)
}

func (u UserServiceClient) Name() string {
	return "user_service"
}

type GetUserReq struct {
	Id int
}

type GetUserResp struct {
	Name string `json:"name"`
}
