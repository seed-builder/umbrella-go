package models

import (
	"github.com/jinzhu/gorm"
	"log"
	"umbrella/utilities"
)

type Site struct {
	gorm.Model
	Base
	Name string
	Province string
	City string
	District string
	Address string
	Longitude string
	Latitude string
	Type int
}

func (m *Site) Save() error{
	utilities.MyDB.Save(m)
	return nil
}

func (m *Site) Remove() error{
	utilities.MyDB.Delete(m)
	return nil
}


func (Site) TableName() string {
	return "sites"
}

func (m *Site) BeforeSave() (err error) {
	log.Println("before save")
	return nil
}

func (m *Site) BeforeDelete() (err error) {
	return nil
}