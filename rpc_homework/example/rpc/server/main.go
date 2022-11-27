package main

import (
	"rpc"
	"rpc/compress/gzip"
	"rpc/serialize/json"
	"rpc/serialize/proto"
)

func main() {
	svr := rpc.NewServer()
	svr.RegisterService(&UserService{})
	svr.RegisterService(&UserServiceProto{})
	svr.RegisterSerializer(json.Serializer{})
	svr.RegisterSerializer(proto.Serializer{})
	svr.RegisterCompresser(&gzip.Compresser{})
	if err := svr.Start(":8081"); err != nil {
		panic(err)
	}
}
