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
	HireDayCash float64
	HireFreeDays uint
	HireExpireDays uint
	Begin time.Time
	End time.Time
	IsDefault uint
	Status uint
}

func (Price) TableName() string {
	return "prices"
}

func (m *Price) Query() *gorm.DB{
	return utilities.MyDB.Model(&Price{})
}