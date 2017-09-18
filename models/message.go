package models

import (
	"github.com/jinzhu/gorm"
	"umbrella/utilities"
	"fmt"
)

type Message struct {
	gorm.Model
	Base
	Category uint
	Level uint
	SiteId uint
	EquipmentId uint
	Channel uint
	Title string
	Content string
	Read uint
}

func (Message) TableName() string {
	return "messages"
}

func (m *Message) Query() *gorm.DB {
	return utilities.MyDB.Model(m)
}

func (m *Message) AddChannelError(sn string, equipment_id uint, site_id uint, channel uint){
	msg := &Message{}
	msg.EquipmentId = equipment_id
	msg.SiteId = site_id
	msg.Channel = channel
	msg.Category = 1
	msg.Level = 2 // warning
	msg.Title = fmt.Sprintf("设备【%s】通道【%d】异常", sn, channel)
	msg.Content = msg.Title
	utilities.MyDB.Create(msg)
}

func (m *Message) AddEquipmentError(sn string, equipment_id uint, site_id uint, content string){
	msg := &Message{}
	msg.EquipmentId = equipment_id
	msg.SiteId = site_id
	msg.Category = 1
	msg.Level = 2 // warning
	msg.Title = fmt.Sprintf("设备【%s】异常： %s", sn, content)
	msg.Content = content
	utilities.MyDB.Create(msg)
}
