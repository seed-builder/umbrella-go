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
	return utilities.MyDB.Model(m)
}

func (m *Equipment) InitChannel() {
	m.ChannelCache = make(map[uint8]uint8, m.Channels)
	umbrella := Umbrella{}
	var have int32
	for i := uint8(1); i <= m.Channels; i ++ {
		var count uint8
		umbrella.Query().Where("status=2 and equipment_id = ? and equipment_channel_num = ?", m.ID, i).Count(&count)
		m.ChannelCache[i] = count
		have =  have + int32(count)
	}
	m.Have = have
	utilities.MyDB.Model(m).Update("have", have)
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
	if len == 0 && m.UsedChannelNum > 0 {
		channelNum = m.UsedChannelNum
	}
	return channelNum
}

func (m *Equipment) InChannel(channelNum uint8){
	n := m.ChannelCache[channelNum]
	m.ChannelCache[channelNum] = n + 1
	m.UsedChannelNum = channelNum
	m.Have = m.Have + 1
	utilities.MyDB.Model(m).Update("have", m.Have )
}

func (m *Equipment) OutChannel(channelNum uint8){
	n := m.ChannelCache[channelNum]
	if n > 0 {
		m.ChannelCache[channelNum] = n - 1
		m.Have = m.Have - 1
		utilities.MyDB.Model(m).Update("have", m.Have)
	}
}

func (m *Equipment) Online(ip string){
	m.Status = EquipmentStatusOnline
	utilities.MyDB.Model(m).Updates(map[string]interface{}{"status":m.Status, "ip": ip})
}

func (m *Equipment) Offline(){
	m.Status = EquipmentStatusOffline
	utilities.MyDB.Model(m).Update("status", m.Status)
}
