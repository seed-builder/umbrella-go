package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type Base struct {
	CreatorId int
	ModifierId int
}

func (m *Base) Query() *gorm.DB{
	return utilities.MyDB
}