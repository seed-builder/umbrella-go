package models

import (
	"github.com/jinzhu/gorm"
	"time"
	"umbrella/utilities"
	"errors"
	"math"
	"net/http"
	"fmt"
	"io/ioutil"
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
	ExpireHours uint
	ExpiredAt time.Time
	HireHours float64
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
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
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
		hire.CalculateFee()
		if status == UmbrellaStatusIn{
			go hire.Notice()
		}
	}
	return status
}

//CalculateHireFee 冻结押金费用
func (m *CustomerHire) FreezeDepositFee(priceId uint){
	price := &Price{}
	price.InitDb(m.db)
	if priceId == 0 {
		price.Query().First(price, "is_default = ?", 1)
	}else{
		price.Query().First(price, priceId)
	}
	if price.ID > 0 {
		m.DepositAmt = price.DepositCash
		m.ExpireHours = price.HireExpireHours
		m.ExpiredAt =  time.Now().Local().Add(time.Hour * time.Duration(price.HireExpireHours))
		m.PriceId = price.ID

		account := &CustomerAccount{}
		account.InitDb(m.db)
		account.Query().First(account, "customer_id = ?", m.CustomerId)
		//utilities.SysLog.Infof("用户【%d】的账户id信息【%d】", m.CustomerId, account.ID)
		if account.ID > 0 {
			res := account.FreezingDeposit(price.DepositCash)
			if res {
				utilities.SysLog.Noticef("成功冻结用户账户【%d】资金【%.2f】", account.ID, price.DepositCash)
				pay := &CustomerPayment{}
				pay.InitDb(m.db)
				pay.PayDeposit(m.ID, m.CustomerId, account.ID, price.DepositCash)
			}else{
				utilities.SysLog.Noticef("冻结账户【%d】资金【%.2f】失败， 余额不足！", m.CustomerId, account.ID, price.DepositCash)
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
			if hours < 0 {
				hours = 0
			}
			var fee float64
			if hours > float64(price.HireFreeHours) {
				calHours := hours - float64(price.HireFreeHours)
				fee = math.Ceil(calHours/price.HireUnitHours) * price.HirePrice
			} else {
				fee = 0
			}
			m.HireAmt = utilities.Round(fee, 2)
			m.HireHours = hours
			m.Query().Updates(map[string]interface{}{
				"hire_hours": hours,
				"hire_amt":   m.HireAmt,
			})
			utilities.SysLog.Noticef("计算费用：租借单【%d】,共【%.1f】小时,租金费用【%.2f】", m.ID, hours, fee)
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
			utilities.SysLog.Noticef("用户【%d】余额【%f】不足，无法结算租借单【%d】租金费用【%.2f】", m.CustomerId, account.BalanceAmt, m.ID, m.HireAmt)
		}
	}
	if complete == 1 {
		m.Query().Updates(map[string]interface{}{
			"status": UmbrellaHireStatusCompleted,
		})
		utilities.SysLog.Noticef("用户【%d】租借单【%d】租金费用【%.2f】结算完成", m.CustomerId, m.ID, m.HireAmt)
		account.ReturnDeposit(m.DepositAmt)
		pay.ReturnDeposit(m.ID, m.CustomerId, account.ID, m.DepositAmt)
	}
}

//发送消息通知已经还伞
func (m *CustomerHire) Notice(){
	// /mobile/customer-hire/return-wechat-send/{id}
	url := utilities.SysConfig.NoticeHost + fmt.Sprintf("/mobile/customer-hire/return-wechat-send/%d", m.ID)
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	utilities.SysLog.Noticef("发送已还伞通知， url【%s】 resp【%s】", url, body)
}

//创建租借记录
func CreateCustomerHire(equipment *Equipment, umbrella *Umbrella, customerId uint, db *gorm.DB) (*CustomerHire, error) {
	if umbrella.ID > 0 {
		entity := &CustomerHire{}
		entity.InitDb(db)
		entity.Entity = entity
		entity.CustomerId = customerId
		entity.CreatorId = customerId
		entity.HireAt = time.Now().Local()//.Add(1 * time.Minute)
		entity.UmbrellaId = umbrella.ID
		entity.HireEquipmentId = equipment.ID
		entity.HireSiteId = equipment.SiteId
		entity.Status = UmbrellaHireStatusNormal
		entity.FreezeDepositFee(equipment.PriceId)
		entity.Save()
		return entity, nil
	}else{
		return  nil, errors.New("伞ID不能为0, 创建租借单失败")
	}
	//go m.FreezeDepositFee(umbrella)
}
