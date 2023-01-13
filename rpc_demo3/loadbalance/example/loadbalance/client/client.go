package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"log"
	"rpc"
	"rpc/loadbalance/example/loadbalance/proto/gen"
	"rpc/loadbalance/roundrobin"
	"rpc/registry/etcd"
)

func main() {
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}
	r, err := etcd.NewRegistry(etcdClient)
	if err != nil {
		panic(err)
	}
	etcdResolver := rpc.NewResolverBuilder(r)

	pickerBudiler := &roundrobin.WeightPickBuilder{}
	builder := base.NewBalancerBuilder("ROUND_ROBIN", pickerBudiler, base.Config{HealthCheck: true})

	balancer.Register(builder)
	cc, err := grpc.Dial("register:///user-service",
		grpc.WithInsecure(),
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"LoadBalancingPolicy": "ROUND_ROBIN"}`))
	if err != nil {
		panic(err)
	}
	client := gen.NewUserServiceClient(cc)
	for i := 0; i < 100; i++ {
		ctx := context.WithValue(context.Background(), "group", "b")
		resp, err := client.GetById(ctx, &gen.GetByIdReq{})
		if err != nil {
			panic(err)
		}
		log.Println(resp.User)
	}

}
