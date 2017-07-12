package utilities

const (
	//失败
	RspStatusFail uint8 = iota
	//成功
	RspStatusSuccess
	//伞过期
	RspStatusUmbrellaExpired
	//非法伞，无法识别
	RspStatusUmbrellaIllegal
	//未登陆授权
	RspStatusUnauth
)
