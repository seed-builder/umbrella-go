package models

import (
	"github.com/jinzhu/gorm"
	"time"
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