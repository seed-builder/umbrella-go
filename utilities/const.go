package utilities

const (
	RspStatusUnknown uint8 = iota
	//成功
	RspStatusSuccess
	//失败
	RspStatusFail
	//伞过期
	RspStatusUmbrellaExpired
	//非法伞，无法识别
	RspStatusUmbrellaIllegal
	//未登陆授权
	RspStatusUnauth
)
