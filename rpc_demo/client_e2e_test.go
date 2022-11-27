package rpc_demo

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	client, err := NewClient(":8081")
	require.NoError(t, err)
	us := &UserServiceClient{}
	err = InitClientProxy(us, client)
	require.NoError(t, err)
	resp, err := us.GetUserId(context.Background(), &GetUserReq{
		Id: 16,
	})
	require.NoError(t, err)
	log.Println(resp)
}
