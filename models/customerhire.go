package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

//(1-租借中, 2-待支付租金 3-已完成, 4-逾期未归还 )
const(
	UmbrellaHireStatusUnknown int32 = iota
	UmbrellaHireStatusNormal
	UmbrellaHireStatusPaying
	UmbrellaHireStatusCompleted
	UmbrellaHireStatusExpired
)

type CustomerHire struct {
	gorm.Model
	Base
	CustomerId uint
	UmbrellaId uint
	HireEquipmentId uint
	HireSiteId uint
	HireAt time.Time
	DepositAmt float64
	ReturnEquipmentId uint
	ReturnSiteId uint
	ReturnAt time.Time
	ExpireDay int32
	ExpiredAt time.Time
	HireDay int32
	HireAmt float64
	Status int32

}


func (CustomerHire) TableName() string {
	return "customer_hires"
}