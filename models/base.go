package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type Base struct {
	CreatorId int
	ModifierId int
	Entity interface{} `gorm:"-"`
}

func (m *Base) Query() *gorm.DB{
	return utilities.MyDB
}

func (m *Base) Save() error{
	utilities.MyDB.Save(m.Entity)
	return nil
}

func (m *Base) Remove() error{
	utilities.MyDB.Delete(m.Entity)
	return nil
}
