package models

import "github.com/jinzhu/gorm"

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
