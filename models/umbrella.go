package models

import(
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
	"umbrella/utilities"
	"strings"
)

// 1-未发放 2-待借中 3-借出中 4-失效（超过还伞时间） 5-异常
const(
	UmbrellaStatusInit int32 = iota + 1
	UmbrellaStatusIn
	UmbrellaStatusOut
	UmbrellaStatusExpired
	UmbrellaStatusAbnormal
)

type Umbrella struct {
	gorm.Model
	Base
	Sn string
	EquipmentId uint
	EquipmentChannelNum uint8
	SiteId uint
	Status int32
	Name string
	Color string
	Logo string
	PriceId uint
}

func NewUmbrella(data map[string]interface{}) *Umbrella  {
	eq := Umbrella{}
	mapstructure.Decode(data, &eq)
	eq.Entity = &eq
	return &eq
}

func (Umbrella) TableName() string {
	return "umbrellas"
}

func (m *Umbrella) BeforeSave() (err error) {
	return nil
}

func (m *Umbrella) BeforeDelete() (err error) {
	return nil
}

func (m *Umbrella) Query() *gorm.DB{
	return utilities.MyDB.Model(&Umbrella{})
}

func (m *Umbrella) Check(umbrellaSn string) uint8 {
	m.Query().First(m, "sn = ?", strings.ToUpper(umbrellaSn))
	if m.ID == 0 {
		return utilities.RspStatusUmbrellaIllegal
	}
	if m.Status == UmbrellaStatusExpired {
		return utilities.RspStatusUmbrellaExpired
	}
	return utilities.RspStatusSuccess
}

//InEquipment 进入设备, 还伞
func (m *Umbrella) InEquipment(equipment *Equipment, umbrellaSn string, channelNum uint8)  uint8 {

	m.Query().First(m, "sn = ?", strings.ToUpper(umbrellaSn))
	if m.ID == 0 {
		return utilities.RspStatusUmbrellaIllegal
	}
	if m.Status == UmbrellaStatusExpired {
		return utilities.RspStatusUmbrellaExpired
	}
	m.Entity = m
	m.EquipmentId = equipment.ID
	m.EquipmentChannelNum = channelNum
	//m.Status = UmbrellaStatusIn

	if m.Status == UmbrellaStatusOut {
		hire := CustomerHire{}
		hire.Entity = &hire
		hire.Query().First(&hire, "umbrella_id = ? and status=1", m.ID)
		now := time.Now()
		if hire.ID > 0 {
			if hire.ExpiredAt.Before(now) {
				hire.Status = UmbrellaHireStatusExpired
				m.Status = UmbrellaStatusExpired
			} else {
				hire.Status = UmbrellaHireStatusPaying
				hire.ReturnAt = time.Now()
				hire.ReturnEquipmentId = equipment.ID
				hire.ReturnSiteId = equipment.SiteId
			}
			hire.Save()
		}else{
			m.Status = UmbrellaStatusIn //UmbrellaStatusAbnormal
		}
	}else if m.Status == UmbrellaStatusInit {
		m.Status = UmbrellaStatusIn
	}
	m.Save()

	if m.Status == UmbrellaStatusIn {
		equipment.InChannel(channelNum)
		return utilities.RspStatusSuccess
	}else{
		return utilities.RspStatusUmbrellaExpired
	}
}

//OutEquipment 出设备
func (m *Umbrella) OutEquipment(equipment *Equipment, umbrellaSn string, channelNum uint8) uint8 {
	m.Query().First(m, "sn = ?", strings.ToUpper(umbrellaSn))
	if m.ID == 0 {
		return utilities.RspStatusUmbrellaIllegal
	}
	m.Entity = m
	m.Status = UmbrellaStatusOut
	m.EquipmentId = equipment.ID
	m.EquipmentChannelNum = channelNum
	m.SiteId = equipment.SiteId
	m.Save()
	return utilities.RspStatusSuccess
}
