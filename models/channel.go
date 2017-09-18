package models

type Channel struct {
	//通道编号
	Id uint8
	//状态
	Status uint8
	//有伞数
	Umbrellas uint8
	//是否有效
	Valid bool
	//发送救援命令次数， 超出3次则记录异常
	RescueTimes int
}
