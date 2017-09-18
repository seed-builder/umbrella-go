package network

//CmdChannelRescueReqPkt 设备通道救援（恢复）命令
type CmdChannelRescueReqPkt struct{
	//session info
	CmdData
}

type CmdChannelRescueRspPkt struct{
	CmdData
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdChannelRescueReqPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_RESCUE
	return p.ToBytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdChannelRescueReqPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdChannelRescueRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_RESCUE_RESP
	return p.ToBytes(p.Status)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdChannelRescueRspPkt) Unpack(data []byte) error {
	if l := len(data); l == 1 {
		p.Status = data[0]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}