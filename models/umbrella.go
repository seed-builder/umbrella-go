package models

import "github.com/jinzhu/gorm"

type Umbrella struct {
	gorm.Model
	Base
	Sn string
	EquipmentId int
	SiteId int
	Status int
	Name string
	Color string
	Logo string

}
