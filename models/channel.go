package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type Channel struct {
	gorm.Model
	Base
	EquipmentId uint
	//通道编号
	Num uint8
	//通道锁状态
	LockStatus uint8
	//有伞数
	Umbrellas uint8
	//是否有效
	Valid bool
	//发送救援命令次数， 超出3次则记录异常
	RescueTimes int
}

func (Channel) TableName() string {
	return "equipment_channels"
}

func (m *Channel) Query() *gorm.DB{
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}

func (m *Channel) UpdateUmbrellas()  {
	utilities.MyDB.Model(m).Update("umbrellas", m.Umbrellas)
}

func (m *Channel) UpdateInfo()  {
	utilities.MyDB.Model(m).Updates(map[string]interface{} {"lock_status": m.LockStatus, "rescue_times": m.RescueTimes, "valid": m.Valid})
}
