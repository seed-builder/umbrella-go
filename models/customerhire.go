package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
)

//(1-初始(未支付押金)，2-租借中, 3-还伞完毕，待支付租金 4-已完成, 5-逾期未归还)
const(
	UmbrellaHireStatusUnknown int32 = iota
	UmbrellaHireStatusInit
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

func (m *CustomerHire) Create(equipment *Equipment, umbrella *Umbrella, customerId uint) error {
	m.Entity = m
	m.CustomerId = customerId
	m.CreatorId = customerId
	m.HireAt = time.Now()
	m.UmbrellaId = m.ID
	m.HireEquipmentId = equipment.ID
	m.HireSiteId = equipment.SiteId
	m.Status = UmbrellaHireStatusInit
	return m.Save()
}