package models

import(
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
	"github.com/mitchellh/mapstructure"
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

func NewUmbrella(data map[string]interface{}) *Umbrella  {
	eq := Umbrella{}
	mapstructure.Decode(data, &eq)
	eq.Entity = &eq
	return &eq
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