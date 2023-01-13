package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"golang.org/x/sync/errgroup"
	"rpc"
	"rpc/loadbalance/example/loadbalance/proto/gen"
	"rpc/registry/etcd"
	"strconv"
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
	var eg errgroup.Group

	for i := 0; i < 3; i++ {
		idx := i
		var group string
		if idx%2 == 0 {
			group = "a"
		} else {
			group = "b"
		}

		eg.Go(func() error {
			xx := uint32(idx + 1)
			server := rpc.NewServer("user-service",
				rpc.ServerWithWeight(xx),
				rpc.ServerWithGroup(group),
				rpc.ServerWithRegistry(r))
			defer server.Close()

			us := &UserService{
				name: fmt.Sprintf("server-%d", idx),
			}
			gen.RegisterUserServiceServer(server, us)
			fmt.Println("启动服务器: " + us.name)
			fmt.Printf("权重: %d", xx)
			// 端口分别是 8081, 8082, 8083
			return server.Start(":" + strconv.Itoa(8081+idx))
		})
	}
	// 正常或者异常退出都会返回
	err = eg.Wait()
	fmt.Println(err)
}
