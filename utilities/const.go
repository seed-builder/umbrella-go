package utilities

const (
	//失败
	RspStatusFail uint8 = iota
	//成功
	RspStatusSuccess
	//超时
	RspStatusTimeout
	//错误
	RspStatusError
	//网络出错
	RspStatusGprs
	//未登陆
	RspStatusNeedAuth
	//伞过期
	RspStatusUmbrellaExpired
	//非法伞，无法识别
	RspStatusUmbrellaIllegal
	//非法设备号
	RspStatusEquipmentSnIllegal
	//非法命令
	RspStatusCmdIllegal
	//未知错误
	RspStatusUnknowError


)

func RspStatusDesc(status uint8) string{
	var desc string
	switch status {
	case RspStatusFail:
		desc = "失败"
	case RspStatusSuccess:
		desc = "成功"
	case RspStatusTimeout:
		desc = "超时"
	case RspStatusError:
		desc = "错误"
	case RspStatusGprs:
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
	}
	return desc
}
