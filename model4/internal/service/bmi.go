package service

import (
	"context"
	"model4/internal/biz"
	"model4/model4/api"
)

type BmiService struct {
	api.UnimplementedBMIServiceServer
	biz *biz.BmiBiz
}
func NewBmiService(biz *biz.BmiBiz)*BmiService{
	return &BmiService{
		biz: biz,
	}
}
func (b *BmiService)BmiInfo(ctx context.Context,request *api.UserInfoRequest)(*api.BMIInfoReply,error){
	bmiinfo,err := b.biz.GetBmiInfoById(request.Uid)
	if err != nil {
		return nil,err
	}
	return &api.BMIInfoReply{
		Bmi: &api.BMI{
			Nickname: bmiinfo.Nickname,
			Bmi: bmiinfo.Bmi,
		},
	},nil
}