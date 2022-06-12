+build wireinject
package main

import (
	"github.com/google/wire"
	"model4/internal/biz"
	"model4/internal/data"
	"model4/internal/service"
)

func InitBmiService()*service.BmiService{
	wire.Build(service.NewBmiService,biz.NewBmiBiz,data.NewBmiRepo,data.NewDB)
	return &service.BmiService{}
}