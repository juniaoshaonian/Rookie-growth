package rpc_demo

import (
	"context"
	"log"
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()
	server.Register(&UserServiceServer{})
	server.Start(":8081")
}

type UserServiceServer struct {
}

func (u *UserServiceServer) Name() string {
	return "user_service"
}

func (u *UserServiceServer) GetUserId(ctx context.Context, req *GetUserReq) (resp *GetUserResp, err error) {
	log.Println(req)
	return &GetUserResp{
		"Tom",
	}, nil
}
