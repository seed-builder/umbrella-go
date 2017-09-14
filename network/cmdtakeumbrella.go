package network

//CmdTakeUmbrellaReqPkt 从通道取伞命令
type CmdTakeUmbrellaReqPkt struct{
	CmdData
}

type CmdTakeUmbrellaRspPkt struct{
	CmdData
	//len 4
	UmbrellaSn []byte
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdTakeUmbrellaReqPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_TAKE_UMBRELLA
	return p.ToBytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdTakeUmbrellaReqPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdTakeUmbrellaRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_TAKE_UMBRELLA_RESP
	var buf []byte
	buf = append(buf, p.UmbrellaSn...)
	buf = append(buf, p.Status)
	return p.ToBytes(buf...)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdTakeUmbrellaRspPkt) Unpack(data []byte) error {
	if l := len(data); l == 5 {
		p.UmbrellaSn = data[:4]
		p.Status = data[4]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}