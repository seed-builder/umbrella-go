package utilities

//状态码
const (
	//成功
	RspStatusSuccess uint8 = 0x00
	//通道忙（有指令未完成
	RspStatusChannelBusy uint8 = 0x01
	//通讯超时
	RspStatusTimeout uint8 = 0x02
	//用户名错误
	RspStatusUserWrong uint8 = 0x03
	//网络出错
	RspStatusGprsErr uint8 = 0x04
	//未登陆
	RspStatusNeedAuth uint8 = 0x05
	//伞过期
	RspStatusUmbrellaExpired uint8 = 0x06
	//非法伞ID，禁止
	RspStatusUmbrellaIllegal uint8 = 0x07
	//非法设备号
	RspStatusEquipmentSnIllegal  uint8 = 0x08
	//非法命令不支持
	RspStatusCmdIllegal uint8 = 0x09
	//未知错误
	RspStatusUnknowError uint8 = 0x0A
	//伞已还（SN检查）
	RspStatusUmbrellaReturned uint8 = 0x0B
	//数据错
	RspStatusDataErr uint8 = 0x0C
	//通道和命令不匹配
	RspStatusNotMatch uint8 = 0x0D

	//通道超时
	RspStatusChannelTimeout uint8 = 0xA0
	//通道锁状态-中间
	RspStatusChannelMiddle uint8 = 0xA1
	//通道锁状态-借伞
	RspStatusChannelBorrow uint8 = 0xA2
	//通道锁状态-还伞
	RspStatusChannelReturn uint8 = 0xA3
	//通道命令不支持
	RspStatusChannelErr uint8 = 0xA4
	//伞SN不匹配
	RspStatusChannelErrSN uint8 = 0xA5
	//通道锁异常
	RspStatusChannelErrLock uint8 = 0xA6


)

func RspStatusDesc(status uint8) string{
	var desc string
	switch status {
	case RspStatusSuccess:
		desc = "成功"
	case RspStatusChannelBusy:
		desc = "通道忙（有指令未完成）"
	case RspStatusTimeout:
		desc = "超时"
	case RspStatusUserWrong:
		desc = "用户名错误"
	case RspStatusGprsErr:
		desc = "网络错误"
	case RspStatusNeedAuth:
		desc = "未登陆"
	case RspStatusUmbrellaExpired:
		desc = "伞过期"
	case RspStatusUmbrellaIllegal:
		desc = "非法伞，无法识别"
	case RspStatusEquipmentSnIllegal:
		desc = "非法设备号"
	case RspStatusCmdIllegal:
		desc = "RspStatusCmdIllegal"
	case RspStatusUnknowError:
		desc = "未知错误"
	case RspStatusUmbrellaReturned:
		desc = "伞已还（SN检查）"
	case RspStatusDataErr:
		desc = "数据错"
	case RspStatusChannelTimeout:
		desc = "通道超时"
	case RspStatusChannelMiddle:
		desc = "通道锁状态-中间"
	case RspStatusChannelBorrow:
		desc = "通道锁状态-借伞"
	case RspStatusChannelReturn:
		desc = "通道锁状态-还伞"
	case RspStatusChannelErr:
		desc = "通道命令不支持"
	case RspStatusNotMatch:
		desc = "通道和命令不匹配"
	case RspStatusChannelErrSN:
		desc = "通道伞SN不匹配"
	case RspStatusChannelErrLock:
		desc = "通道锁异常"
	}
	return desc
}
