package models

import (
	"github.com/jinzhu/gorm"
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
