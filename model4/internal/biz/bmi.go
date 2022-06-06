package biz

import (
	"github.com/pkg/errors"
	"model4/internal/data"
)
type BmiBiz struct {
	repo *data.BmiRepo
}
func NewBmiBiz(repo *data.BmiRepo)*BmiBiz {
	return &BmiBiz{
		repo: repo,
	}
}
func (bb *BmiBiz)GetBmiInfoById(id uint64)(*BmiDO,error){
	if id == 0 {
		return nil,errors.New("invalid id")
	}
	u,err := bb.repo.GerUser(id)
	if err != nil {
		return nil,err
	}
	return &BmiDO{
		Nickname: u.Nickname,
		Bmi: u.Bmi,
	},nil
}

type BmiDO struct {
	Nickname string
	Bmi uint64
}
type MyService interface {

}

type Option func(db *service)MyService

type service struct {

}

func NewDB(opts... Option)MyService{
	return &service{}
}



