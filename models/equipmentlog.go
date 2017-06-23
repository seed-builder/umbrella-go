package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type EquipmentLog struct {
	gorm.Model
	Base
	EquipmentId int
	SiteId int
	ApiName string
	Code string
	Type string
	Content string
}

func (m *EquipmentLog) Save() error{
	utilities.MyDB.Save(m)
	return nil
}

func (m *EquipmentLog) Remove() error{
	utilities.MyDB.Delete(m)
	return nil
}

func (EquipmentLog) TableName() string {
	return "equipment_logs"
}

func (m *EquipmentLog) BeforeSave() (err error) {
	m.Entity = m
	return nil
}

func (m *EquipmentLog) BeforeDelete() (err error) {
	m.Entity = m
	return nil
}