package models

import (
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"umbrella/utilities"
)

const (
	EquipmentStatusNone int32 = iota + 1
	EquipmentStatusUse
	EquipmentStatusOnline
	EquipmentStatusOffline
	EquipmentStatusBug
)

type Equipment struct {
	gorm.Model
	Base
	Sn string
	SiteId uint
	Capacity int32
	Have int32
	Type int32
	Ip string
	Status int32
	Channels uint8
	ChannelCache map[uint8]uint8 `gorm:"-"`
	UsedChannelNum uint8 `gorm:"-"`
}

func NewEquipment(data map[string]interface{}) *Equipment  {
	eq := Equipment{}
	mapstructure.Decode(data, &eq)
	eq.Entity = &eq
	return &eq
}

func (Equipment) TableName() string {
	return "equipments"
}

func (m *Equipment) BeforeSave() (err error) {
	return nil
}

func (m *Equipment) BeforeDelete() (err error) {
	return nil
}

func (m *Equipment) Query() *gorm.DB{
	return utilities.MyDB.Model(&Equipment{})
}

func (m *Equipment) InitChannel() {
	m.ChannelCache = make(map[uint8]uint8, m.Channels)
	umbrella := Umbrella{}
	for i := uint8(1); i <= m.Channels; i ++ {
		var count uint8
		umbrella.Query().Where("equipment_id = ? and equipment_channel_num = ?", m.ID, i).Count(&count)
		m.ChannelCache[i] = count
	}
}

//ChooseChannel 选择伞保有量最多的通道
func (m *Equipment) ChooseChannel() uint8 {
	var len uint8
	channelNum := uint8(1)
	for n, l :=  range m.ChannelCache {
		if n != m.UsedChannelNum && l > len {
			channelNum = n
			len = l
		}
	}
	return channelNum
}

func (m *Equipment) InChannel(channelNum uint8){
	n := m.ChannelCache[channelNum]
	m.ChannelCache[channelNum] = n + 1
	m.UsedChannelNum = channelNum
}

func (m *Equipment) OutChannel(channelNum uint8){
	n := m.ChannelCache[channelNum]
	m.ChannelCache[channelNum] = n - 1
}


func (m *Equipment) Online(){
	m.Status = EquipmentStatusOnline
	utilities.MyDB.Model(m).Update("status", m.Status)
}

func (m *Equipment) Offline(){
	m.Status = EquipmentStatusOffline
	utilities.MyDB.Model(m).Update("status", m.Status)
}
