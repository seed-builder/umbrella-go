package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

type Base struct {
	CreatorId uint
	ModifierId uint
	Entity interface{} `gorm:"-"`
	db *gorm.DB `gorm:"-"`
}

func (m *Base) InitDb(db *gorm.DB) {
	m.db = db
}

func (m *Base) Query() *gorm.DB{
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db
}

func (m *Base) Save() error{
	m.Query().Save(m.Entity)
	return nil
}

func (m *Base) Remove() error{
	m.Query().Delete(m.Entity)
	return nil
}

func (m *Base) Update(properties map[string]interface{}){
	m.Query().Model(m.Entity).Updates(properties)
}
