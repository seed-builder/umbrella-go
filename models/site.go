package models

import (
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
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
	Type int32
}

func NewSite(data map[string]interface{}) *Site  {
	site := Site{}
	mapstructure.Decode(data, &site)
	site.Entity = &site
	return &site
}

func (m Site) TableName() string {
	return "sites"
}

func (m *Site) BeforeSave() (err error) {
	return nil
}

func (m *Site) BeforeDelete() (err error) {
	return nil
}

func (m *Site) Query() *gorm.DB{
	return utilities.MyDB.Model(&Site{})
}