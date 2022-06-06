package data

import (
	"github.com/gotomicro/ego-component/egorm"
)
type BmiRepo struct {
	db *egorm.Component
}

func NewBmiRepo(db *egorm.Component)*BmiRepo{
	return &BmiRepo{
		db: db,
	}
}
func (bb *BmiRepo)GerUser(id uint64)(*BmiPO,error){
	return &BmiPO{
		Nickname: "ZZ",
		Bmi: uint64(19),
	},nil

}
type BmiPO struct {
	Nickname string
	Bmi uint64

}
type BmiDO struct {
	Nickname string
	Bmi  uint64
}