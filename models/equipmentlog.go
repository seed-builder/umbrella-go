package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
	"github.com/mitchellh/mapstructure"
	"regexp"
	"fmt"
)

type EquipmentLog struct {
	gorm.Model
	Base
	Level int
	Content string
	EquipmentSn string
}

func NewEquipmentLog(data map[string]interface{}) *EquipmentLog  {
	eq := EquipmentLog{}
	mapstructure.Decode(data, &eq)
	eq.Entity = &eq
	return &eq
}

func (EquipmentLog) TableName() string {
	return "equipment_logs"
}

func (m *EquipmentLog) BeforeSave() (err error) {
	return nil
}

func (m *EquipmentLog) BeforeDelete() (err error) {
	return nil
}

func (m *EquipmentLog) Query() *gorm.DB{
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}

func (m *EquipmentLog) NewLog(level int, content string) bool {
	log := EquipmentLog{ Level: level, Content: content}
	if sn, ok := ParseEquipmentSn(content); ok {
		log.EquipmentSn = sn
	}
	go utilities.MyDB.Create(&log)
	return true
}

func ParseEquipmentSn(content string) (string, bool) {
	reg := regexp.MustCompile(`设备【(\w+)】`)
	res := reg.FindSubmatch([]byte(content))
	if l := len(res); l == 2 {
		sn := fmt.Sprintf("%q", res[1])
		return sn, true
	}
	return "", false
}