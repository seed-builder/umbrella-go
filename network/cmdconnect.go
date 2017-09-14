package network

import (
	"errors"
	"umbrella/utilities"
)

const(
	EquipmentSnLen int = 11
)

//CmdConnectReqPkt is the connect request packet
type CmdConnectReqPkt struct {
	//session info
	CmdData
	EquipmentSn string
}

type CmdConnectRspPkt struct {
	//session info
	CmdData
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdConnectReqPkt) Pack(seqId uint8) ([]byte, error) {
	utilities.SysLog.Infof("connect pack seqid %d", seqId)
	p.SeqId = seqId
	p.CmdId = CMD_CONNECT
	body := []byte(p.EquipmentSn)
	return p.ToBytes(body...)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdConnectReqPkt) Unpack(data []byte) error {
	if l := len(data); l == EquipmentSnLen {
		p.EquipmentSn = string(data[:])
		return nil
	}else{
		return errors.New("Equipment Sn Length is wrong!")
	}
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdConnectRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CONNECT_RESP
	return p.ToBytes(p.Status)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdConnectRspPkt) Unpack(data []byte) error {
	if l := len(data); l == 1 {
		p.Status = data[0]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}
