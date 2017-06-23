package models

import "github.com/jinzhu/gorm"

type IDB interface {
	Save() error
	Remove() error
	Query() *gorm.DB
}