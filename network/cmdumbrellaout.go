package network

const (
	UmbrellaSnLen int = 4
)

type CmdUmbrellaOutReqPkt struct{
	CmdData
	//EquipmentSn string
	//len 4
	UmbrellaSn []byte
}

type CmdUmbrellaOutRspPkt struct{
	CmdData
	//len 4
	UmbrellaSn []byte
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutReqPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_UMBRELLA_OUT
	return p.ToBytes(p.UmbrellaSn...)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaOutReqPkt) Unpack(data []byte) error {
	if l := len(data); l == 4 {
		p.UmbrellaSn = data[:]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_UMBRELLA_OUT_RESP
	var buf []byte
	buf = append(buf, p.UmbrellaSn...)
	buf = append(buf, p.Status)
	return p.ToBytes(buf...)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaOutRspPkt) Unpack(data []byte) error {
	l := len(data)
	if l == 5 {
		p.UmbrellaSn = data[:4]
		p.Status = data[4]
	}else if l == 1{
		p.Status = data[0]
	}else{
		return ErrCmdDataLengthWrong
	}
	return nil
}