package client

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"rpc"
	"rpc_demo2/example/proto/gen"
)

func main() {

	rsBuilder := rpc.NewResolverBuilder()
	cc, err := grpc.Dial("register:///user-service", grpc.WithResolvers(rsBuilder), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	usClient := gen.NewUserServiceClient(cc)
	resp, err := usClient.GetById(context.Background(), &gen.GetByIdReq{
		Id: 12,
	})
	if err != nil {
		panic(err)
	}
	log.Println(resp)

}
