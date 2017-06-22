package models

import "github.com/jinzhu/gorm"

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