package models

import (
	"github.com/jinzhu/gorm"
	//"umbrella/utilities"
	"github.com/mitchellh/mapstructure"
)

type Equipment struct {
	gorm.Model
	Base
	Sn string
	SiteId int
	Capacity int
	Have int
	Type int
	Ip string
	Status int
}


func NewEquipment(data map[string]interface{}) *Equipment  {
	eq := Equipment{}
	mapstructure.Decode(data, &eq)
	eq.Entity = &eq
	return &eq
}

func (Equipment) TableName() string {
	return "equipments"
}

func (m *Equipment) BeforeSave() (err error) {
	m.Entity = m
	return nil
}

func (m *Equipment) BeforeDelete() (err error) {
	m.Entity = m
	return nil
}