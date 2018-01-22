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
	ServerHttpBase string
	ChannelCache map[uint8]*Channel `gorm:"-"`
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
	if m.db == nil{
		m.db = utilities.MyDB
	}
	return m.db.Model(m)
}

func (m *Equipment) InitChannel() {
	var channels []Channel
	channel := &Channel{}
	channel.Query().Find(&channels, "equipment_id = ?", m.ID)
	m.ChannelCache = make(map[uint8]*Channel, len(channels))

	umbrella := &Umbrella{}
	var have int32
	for i := 0; i < len(channels); i ++ {
		var count uint8
		c := channels[i]
		umbrella.Query().Where("status=2 and equipment_id = ? and equipment_channel_num = ?", m.ID, c.Num).Count(&count)
		c.Umbrellas = count
		m.ChannelCache[c.Num] = &c //&Channel{ Num: , Umbrellas: count, }
		have =  have + int32(count)
		go c.UpdateUmbrellas()
	}
	m.Have = have
	utilities.MyDB.Model(m).Update("have", have)
}

//ChooseChannel 选择伞保有量最多的通道
func (m *Equipment) ChooseChannel() uint8 {
	var len uint8
	channelNum := uint8(0)
	for n, l :=  range m.ChannelCache {
		if n != m.UsedChannelNum && l.Umbrellas > len && m.CheckIsUseful(n) {
			channelNum = n
			len = l.Umbrellas
		}
	}
	if len == 0 && m.UsedChannelNum > 0 && m.CheckIsUseful(m.UsedChannelNum) {
		channelNum = m.UsedChannelNum
	}
	return channelNum
}

func (m *Equipment) InChannel(channelNum uint8){
	n := m.ChannelCache[channelNum]
	n.Umbrellas = n.Umbrellas + 1
	go n.UpdateUmbrellas()

	m.UsedChannelNum = channelNum
	m.Have = m.Have + 1
	utilities.MyDB.Model(m).Update("have", m.Have )
}

func (m *Equipment) OutChannel(channelNum uint8, db *gorm.DB) error {
	n := m.ChannelCache[channelNum]
	if n.Umbrellas > 0 {
		n.Umbrellas = n.Umbrellas - 1
		go n.UpdateUmbrellas()

		m.Have = m.Have - 1
		if db != nil {
			return db.Model(m).Update("have", m.Have).Error
		}else {
			return utilities.MyDB.Model(m).Update("have", m.Have).Error
		}
	}
	return nil
}

func (m *Equipment) Online(ip string){
	m.Status = EquipmentStatusOnline
	utilities.MyDB.Model(m).Updates(map[string]interface{}{"status":m.Status, "ip": ip, "server_http_base": m.ServerHttpBase})
}

func (m *Equipment) Offline(){
	m.Status = EquipmentStatusOffline
	utilities.MyDB.Model(m).Update("status", m.Status)
}

func (m *Equipment)SetChannelStatus(num uint8, status uint8) bool {
	n, ok := m.ChannelCache[num]
	if ok && n.Valid {
		//rescue := false
		//n.LockStatus = status

		if status == utilities.RspStatusChannelTimeout || status ==  utilities.RspStatusTimeout || status == utilities.RspStatusChannelErrLock {
			n.RescueTimes ++
			if n.RescueTimes >= 3 {
				n.LockStatus = status
			}
		} else {
			//n.Valid = true
			n.LockStatus = status
			n.RescueTimes = 0
		}
		go n.UpdateInfo()
	}
	return false
}

func (m *Equipment)SetChannelValid(num uint8, valid bool) {
	n, ok := m.ChannelCache[num]
	if ok {
		n.LockStatus = utilities.RspStatusChannelBorrow
		n.Valid = valid
	}
}

// 检测是否可用
func (m *Equipment) CheckIsUseful(channelNum uint8) bool {
	n := m.ChannelCache[channelNum]
	return n.Valid && n.LockStatus > utilities.RspStatusChannelTimeout && n.LockStatus < utilities.RspStatusChannelErrLock
}

// 检查是否有效
func (m *Equipment) CheckValid(channelNum uint8) bool {
	n := m.ChannelCache[channelNum]
	return n.Valid
}