package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
	"errors"
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
	HireDay float64
	HireAmt float64
	Status uint
	PriceId uint

	HireEquipment Equipment `gorm:"ForeignKey:HireEquipmentId"`
	ReturnEquipment Equipment `gorm:"ForeignKey:ReturnEquipmentId"`
}


func (CustomerHire) TableName() string {
	return "customer_hires"
}

func (m *CustomerHire) Query() *gorm.DB{
	return utilities.MyDB.Model(m)
}

func (m *CustomerHire) Create(equipment *Equipment, umbrella *Umbrella, customerId uint) (bool, error) {
	if umbrella.ID > 0 {
		m.Entity = m
		m.CustomerId = customerId
		m.CreatorId = customerId
		m.HireAt = time.Now().Local()//.Add(1 * time.Minute)
		m.UmbrellaId = umbrella.ID
		m.HireEquipmentId = equipment.ID
		m.HireSiteId = equipment.SiteId
		m.Status = UmbrellaHireStatusNormal
		m.Save()
		return true, nil
	}else{
		return  false, errors.New("伞ID不能为0")
	}
	//go m.FreezeDepositFee(umbrella)
}

//还伞
func (hire *CustomerHire) UmbrellaReturn(umbrellaId uint, equipmentId uint, siteId uint) int32 {
	hire.Query().First(&hire, "umbrella_id = ? and status = 2", umbrellaId)
	now := time.Now().Local()
	status := UmbrellaStatusIn
	if hire.ID > 0 {
		if hire.ExpiredAt.Before(now) {
			hire.Status = UmbrellaHireStatusExpired
			status = UmbrellaStatusExpired
		} else {
			hire.Status = UmbrellaHireStatusPaying
			hire.ReturnAt = time.Now().Local()
			hire.ReturnEquipmentId = equipmentId
			hire.ReturnSiteId = siteId
			status = UmbrellaStatusIn //UmbrellaStatusAbnormal
		}
		hire.Query().Save(hire)
		go hire.CalculateFee()
	}
	return status
}

//CalculateHireFee 结算租借押金费用
func (m *CustomerHire) FreezeDepositFee(umbrella *Umbrella){
	price := &Price{}
	if umbrella.PriceId == 0 {
		price.Query().First(price, "is_default = ?", 1)
	}else{
		price.Query().First(price, umbrella.PriceId)
	}
	if price.ID > 0 {
		m.DepositAmt = price.DepositCash
		m.ExpireDay = price.HireExpireDays
		m.ExpiredAt =  time.Now().Local().Add(time.Hour * 24 * time.Duration(price.HireExpireDays))
		m.Query().Updates(map[string]interface{}{
			"deposit_amt": m.DepositAmt,
			"expire_day": m.ExpireDay,
			"expired_at": m.ExpiredAt,
			"price_id": price.ID,
		})

		account := &CustomerAccount{}
		account.Query().First(account, "customer_id = ?", m.CustomerId)
		//utilities.SysLog.Infof("用户【%d】的账户id信息【%d】", m.CustomerId, account.ID)
		if account.ID > 0 {
			res := account.FreezingDeposit(price.DepositCash)
			utilities.SysLog.Noticef("用户【%d】的冻结账户【%d】资金【%f】", m.CustomerId, account.ID, price.DepositCash)
			if res {
				pay := &CustomerPayment{}
				pay.PayDeposit(m.ID, m.CustomerId, account.ID, price.DepositCash)
			}
		}
	}
}

//CalculateHireFee 结算租金费用
func (m *CustomerHire) CalculateFee(){
	if m.Status == UmbrellaHireStatusPaying {
		price := &Price{}
		price.Query().First(price, m.PriceId)
		if price.ID > 0 {
			calAt := m.HireAt.Add(time.Duration(price.DelaySeconds) * time.Second)
			hours := m.ReturnAt.Sub(calAt).Hours()
			var fee float64
			var days float64
			if hours > 0 {
				days := utilities.Round(hours/24, 2)
				if days > float64(price.HireFreeDays) {
					fee = (days - float64(price.HireFreeDays)) * price.HireDayCash
				} else {
					days = 0
				}
			} else {
				days = 0
				fee = 0
			}
			m.HireAmt = utilities.Round(fee, 2)
			m.HireDay = days
			m.Query().Updates(map[string]interface{}{
				"hire_day": days,
				"hire_amt": m.HireAmt,
			})
			utilities.SysLog.Noticef("计算费用：租借单【%d】,共【%.1f】天,租金费用【%.2f】", m.ID, days, fee)
			//结算
			m.Settle()

		}
	}
}

//结算， 如果用户有余额，则自动结算，归还押金
func(m *CustomerHire) Settle(){
	complete := 0
	account := &CustomerAccount{}
	account.Query().First(account, "customer_id = ?", m.CustomerId)
	pay := &CustomerPayment{}
	if m.HireAmt == 0 {
		complete = 1
	}else {
		if account.ID > 0 && account.BalanceAmt >= m.HireAmt {
			account.Query().Updates(map[string]interface{}{
				"deposit": gorm.Expr("balance_amt - ?", m.HireAmt),
			})
			pay.PayFee(m.ID, m.CustomerId, account.ID, m.HireAmt)
			complete = 1
		} else {
			utilities.SysLog.Noticef("用户【%d】余额【%f】不足，无法结算租借单【%d】租金费用【%f】", m.CustomerId, account.BalanceAmt, m.ID, m.HireAmt)
		}
	}
	if complete == 1 {
		m.Query().Updates(map[string]interface{}{
			"status": UmbrellaHireStatusCompleted,
		})
		utilities.SysLog.Noticef("用户【%d】租借单【%d】租金费用【%f】结算完成", m.CustomerId, m.ID, m.HireAmt)
		account.ReturnDeposit(m.DepositAmt)
		pay.ReturnDeposit(m.ID, m.CustomerId, account.ID, m.DepositAmt)
	}
}