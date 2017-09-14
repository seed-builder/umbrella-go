package models

type Channel struct {
	//通道编号
	Id uint8
	//状态
	Status uint8
	//有伞数
	Umbrellas uint8
	//超时次数
	Timeouts uint8
}
