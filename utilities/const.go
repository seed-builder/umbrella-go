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

)
