package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
	"time"
	"fmt"
)

type CustomerPayment struct {
	gorm.Model
	Base
	CustomerId uint
	CustomerAccountId uint
	Sn string
	OuterOrderSn string
	PaymentChannel uint //支付渠道 1-微信支付 2-支付宝, 3-余额支付
	Amt float64
	//流水类型 1-充值（收入）， 2-押金充值， 3-押金支出， 4-押金退回， 5-借伞租金支出， 6-账户提现")
	Type uint
	//支付状态（1-未支付, 2-已支付, 3-支付失败）
	Status uint
	Remark string
	ReferenceId uint
	ReferenceType string

}

func (CustomerPayment) TableName() string {
	return "customer_payments"
}

func (m *CustomerPayment) Query() *gorm.DB{
	return utilities.MyDB.Model(m)
}

//支付费用
func (cp *CustomerPayment) PayFee(hireId uint, customerId uint, accountId uint, amt float64) {
	m := &CustomerPayment{}
	m.Entity = m
	m.CustomerId = customerId
	m.CustomerAccountId = accountId
	m.Amt = amt
	m.Type = 5
	m.Status = 2
	m.Remark = "借伞租金支出"
	m.ReferenceId = hireId
	m.ReferenceType = "customer_hire"
	m.PaymentChannel = 3
	m.Sn = m.SN(m.CustomerId, 3)
	m.Save()
}

//支付押金
func (cp *CustomerPayment) PayDeposit(hireId uint, customerId uint, accountId uint, amt float64) {
	m := &CustomerPayment{}
	m.Entity = m
	m.CustomerId = customerId
	m.CustomerAccountId = accountId
	m.Amt = amt
	m.Type = 3
	m.Status = 2
	m.Remark = "押金支出"
	m.ReferenceId = hireId
	m.ReferenceType = "customer_hire"
	m.Sn = m.SN(m.CustomerId, 1)
	m.Save()
}

//退回押金
func (cp *CustomerPayment) ReturnDeposit(hireId uint, customerId uint, accountId uint, amt float64) {
	m := &CustomerPayment{}
	m.Entity = m
	m.CustomerId = customerId
	m.CustomerAccountId = accountId
	m.Amt = amt
	m.Type = 4
	m.Status = 2
	m.Remark = "押金退回"
	m.ReferenceId = hireId
	m.ReferenceType = "customer_hire"
	m.Sn = m.SN(m.CustomerId, 2)
	m.Save()
}

func (m *CustomerPayment) SN(customerId uint, typ int) string {
	prefix := "YO"
	switch typ {
	case 1:
		prefix = "YO"
	case 2:
		prefix = "YB"
	case 3:
		prefix = "HO"
	}
	return prefix + fmt.Sprintf("%05d", customerId) + time.Now().Format("20060102150405")
}
