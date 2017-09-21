package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
)

type Customer struct {
	gorm.Model
	Base
	Mobile string
	Openid string
	Nickname string
	HeadImgUrl string
	Password string
	LoginTime int32
	Gender int32
	BirthDay time.Time
	Address string
	Remark string
	Country string
	Province string
	City string
}

func (Customer) TableName() string {
	return "customers"
}

func (m *Customer) Query() *gorm.DB{
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}