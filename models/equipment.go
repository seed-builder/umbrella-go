package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
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


func (m *Equipment) Save() error{
	utilities.MyDB.Save(m)
	return nil
}

func (m *Equipment) Remove() error{
	utilities.MyDB.Delete(m)
	return nil
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