package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
	"strings"
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

func (m *Customer) CanBorrowUmbrella(customerId uint, sn string) bool {
	account := &CustomerAccount{}
	account.Query().First(account,"customer_id = ?", customerId)
	umbrella := &Umbrella{}
	umbrella.Query().First(umbrella, "sn = ?", strings.ToUpper(sn))
	price := &Price{}
	if umbrella.PriceId == 0 {
		price.Query().First(price, "is_default = ?", 1)
	}else{
		price.Query().First(price, umbrella.PriceId)
	}
	return price.DepositCash <= account.Deposit
}

func (m *Customer) CanBorrowFromEquipment(customerId uint, priceId uint) (bool, float64) {
	account := &CustomerAccount{}
	account.Query().First(account,"customer_id = ?", customerId)
	price := &Price{}
	if priceId == 0 {
		price.Query().First(price, "is_default = ?", 1)
	}else{
		price.Query().First(price, priceId)
	}
	return price.DepositCash <= account.Deposit, price.DepositCash
}