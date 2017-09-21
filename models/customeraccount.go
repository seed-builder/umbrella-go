package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type CustomerAccount struct {
	gorm.Model
	Base
	Sn string
	CustomerId uint
	//可用余额
	BalanceAmt float64
	//押金余额
	Deposit float64
	//冻结押金
	FreezeDeposit float64
}

func (CustomerAccount) TableName() string {
	return "customer_accounts"
}

func (m *CustomerAccount) Query() *gorm.DB{
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}

//冻结押金
func (m *CustomerAccount) FreezingDeposit(amt float64) bool {
	if m.Deposit >= amt {
		m.Query().Updates(map[string]interface{}{
			"deposit": gorm.Expr("deposit - ?", amt),
			"freeze_deposit": gorm.Expr("freeze_deposit + ?", amt),
		})
		return true
	}
	return false
}

//退回押金
func (m *CustomerAccount) ReturnDeposit(amt float64) {
	m.Query().Updates(map[string]interface{}{
		"deposit": gorm.Expr("deposit + ?", amt),
		"freeze_deposit": gorm.Expr("freeze_deposit - ?", amt),
	})
	utilities.SysLog.Noticef("用户【%d】退回押金【%f】", m.CustomerId, amt)
}
