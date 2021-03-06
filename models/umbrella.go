package models

import(
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
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
	Number string
	Sn string
	EquipmentId uint
	EquipmentChannelNum uint8
	SiteId uint
	Status int32
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
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}

func (m *Umbrella) Check(umbrellaSn string) uint8 {
	m.Query().First(m, "sn = ?", strings.ToUpper(umbrellaSn))
	if m.ID == 0 {
		return utilities.RspStatusUmbrellaIllegal
	}
	if m.Status == UmbrellaStatusExpired {
		return utilities.RspStatusUmbrellaExpired
	}
	if m.Status == UmbrellaStatusIn {
		return utilities.RspStatusSuccess
	}
	if m.Status == UmbrellaStatusOut {
		return utilities.RspStatusUmbrellaStatusOut
	}
	return utilities.RspStatusUmbrellaStatusErr
}

//InEquipment 进入设备, 还伞
func (m *Umbrella) InEquipment(equipment *Equipment, umbrellaSn string, channelNum uint8)  uint8 {
	//
	m.Query().First(m, "sn = ?", strings.ToUpper(umbrellaSn))
	if m.ID == 0 {
		utilities.SysLog.Warningf("设备【%s】非法伞编号【%s】,禁止进入通道【%d】", equipment.Sn, umbrellaSn, channelNum)
		return utilities.RspStatusUmbrellaIllegal
	}
	if m.Status == UmbrellaStatusExpired {
		utilities.SysLog.Warningf("设备【%s】伞过期编号【%s】,禁止进入通道【$d】", equipment.Sn, umbrellaSn, channelNum)
		return utilities.RspStatusUmbrellaExpired
	}
	if m.Status == UmbrellaStatusIn {
		if m.EquipmentChannelNum == channelNum {
			return utilities.RspStatusSuccess
		}else{
			equipment.OutChannel(m.EquipmentChannelNum, nil)
		}
	}
	equipment.InChannel(channelNum)
	//
	m.Entity = m
	m.EquipmentId = equipment.ID
	m.SiteId = equipment.SiteId
	m.EquipmentChannelNum = channelNum
	//m.Status = UmbrellaStatusIn

	if m.Status == UmbrellaStatusOut {
		hire := &CustomerHire{}
		m.Status = hire.UmbrellaReturn(m.ID, equipment.ID, equipment.SiteId)
	}else {
		m.Status = UmbrellaStatusIn
	}
	m.Save()

	if m.Status == UmbrellaStatusIn {
		return utilities.RspStatusSuccess
	}else{
		return utilities.RspStatusUmbrellaExpired
	}
}

//OutEquipment 出设备
func (m *Umbrella) OutEquipment(equipment *Equipment, umbrellaSn string, channelNum uint8) uint8 {
	m.Query().First(m, "sn = ?", strings.ToUpper(umbrellaSn))
	if m.ID == 0 {
		utilities.SysLog.Warningf("设备【%s】非法伞编号【%s】,禁止出通道【%d】", equipment.Sn, umbrellaSn, channelNum)
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
