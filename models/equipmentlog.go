package models

import (
	"github.com/jinzhu/gorm"
	//"umbrella/utilities"
	"github.com/mitchellh/mapstructure"
)

type EquipmentLog struct {
	gorm.Model
	Base
	EquipmentId uint
	SiteId uint
	ApiName string
	Code string
	Type string
	Content string
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
	m.Entity = m
	return nil
}

func (m *EquipmentLog) BeforeDelete() (err error) {
	m.Entity = m
	return nil
}