package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
)

//(1-初始(未支付押金)，2-租借中, 3-还伞完毕，待支付租金 4-已完成, 5-逾期未归还)
const(
	UmbrellaHireStatusUnknown uint = iota
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
	ExpireDay uint
	ExpiredAt time.Time
	HireDay uint
	HireAmt float64
	Status uint

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
	m.UmbrellaId = umbrella.ID
	m.HireEquipmentId = equipment.ID
	m.HireSiteId = equipment.SiteId
	m.Status = UmbrellaHireStatusInit
	if umbrella.PriceId == 0 {
		price := &Price{}
		price.Query().First(price, "is_default = ?", 1)
		if price.ID > 0 {
			m.DepositAmt = price.DepositCash
			m.ExpireDay = price.HireExpireDays
			m.ExpiredAt =  time.Now().Add(time.Hour * 24 * time.Duration(price.HireExpireDays))
		}
	}
	return m.Save()
}