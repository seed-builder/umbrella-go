package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type Base struct {
	CreatorId uint
	ModifierId uint
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

func (m *Base) Update(properties map[string]interface{}){
	utilities.MyDB.Model(m.Entity).Updates(properties)
}
