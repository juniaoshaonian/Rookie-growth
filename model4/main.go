package main

import (
	"fmt"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego/server/egrpc"
	"model4/model4/api"
)

func main(){
	if err := ego.New().Serve(func()*egrpc.Component{
		srv := egrpc.Load("server.grpc").Build()
		api.RegisterBMIServiceServer(srv,InitBmiService())
		return srv
	}()).Run();err != nil {
		fmt.Println(err)
	}
}