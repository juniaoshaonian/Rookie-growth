package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	rpc_service "rpc"
	"rpc/example/proto/gen"
)

func main() {
	rs := rpc_service.NewResolverBuilder()
	cc, err := grpc.Dial("register:///user-service", grpc.WithResolvers(rs))
	if err != nil {
		panic(err)
	}
	client := gen.NewUserServiceClient(cc)
	resp, err := client.GetById(context.Background(), &gen.GetByIdReq{
		Id: 13,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
