package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
)

//(1-初始，2-未拿伞租借失败，3-租借中, 4-还伞完毕，待支付租金 5-已完成, 6-逾期未归还)
const(
	UmbrellaHireStatusUnknown int32 = iota
	UmbrellaHireStatusInit
	UmbrellaHireStatusFail
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

	HireEquipment Equipment `gorm:"ForeignKey:HireEquipmentId"`
	ReturnEquipment Equipment `gorm:"ForeignKey:ReturnEquipmentId"`

}


func (CustomerHire) TableName() string {
	return "customer_hires"
}

func (m *CustomerHire) Query() *gorm.DB{
	return utilities.MyDB.Model(&CustomerHire{})
}