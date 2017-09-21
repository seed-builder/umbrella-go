package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
)

type Price struct {
	gorm.Model
	Base
	Name string
	DepositCash float64
	//HireDayCash float64
	//HireFreeDays uint
	HireExpireHours uint
	Begin time.Time
	End time.Time
	IsDefault uint
	Status uint
	DelaySeconds uint
	//租金价格
	HirePrice float64
	//价格计算单位（小时）
	HireUnitHours float64
	//免费小时数
	HireFreeHours float64
}

func (Price) TableName() string {
	return "prices"
}

func (m *Price) Query() *gorm.DB{
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}