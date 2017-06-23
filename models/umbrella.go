package models

import(
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
)

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


func (m *Umbrella) Save() error{
	utilities.MyDB.Save(m)
	return nil
}

func (m *Umbrella) Remove() error{
	utilities.MyDB.Delete(m)
	return nil
}

func (Umbrella) TableName() string {
	return "sites"
}

func (m *Umbrella) BeforeSave() (err error) {
	return nil
}

func (m *Umbrella) BeforeDelete() (err error) {
	return nil
}