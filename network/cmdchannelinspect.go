package network

//CmdInspectReqPkt 设备通道检查命令
type CmdChannelInspectReqPkt struct{
	//session info
	CmdData
}

type CmdChannelInspectRspPkt struct{
	CmdData
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdChannelInspectReqPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_INSPECT
	return p.ToBytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdChannelInspectReqPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdChannelInspectRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_INSPECT_RESP
	return p.ToBytes(p.Status)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdChannelInspectRspPkt) Unpack(data []byte) error {
	if l := len(data); l == 1 {
		p.Status = data[0]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}