package network

type CmdUmbrellaInReqPkt struct{
	CmdData
	//4字节
	UmbrellaSn []byte
}

type CmdUmbrellaInRspPkt struct{
	CmdData
	//4字节
	UmbrellaSn []byte
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaInReqPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_UMBRELLA_IN
	return p.ToBytes(p.UmbrellaSn...)

}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaInReqPkt) Unpack(data []byte) error {
	if l := len(data); l == 4 {
		p.UmbrellaSn = data[:]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaInRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_CHANNEL_UMBRELLA_IN_RESP
	var buf []byte
	buf = append(buf, p.UmbrellaSn...)
	buf = append(buf, p.Status)
	return p.ToBytes(buf...)
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaInRspPkt) Unpack(data []byte) error {
	if l := len(data); l == 5 {
		p.UmbrellaSn = data[:4]
		p.Status = data[4]
		return nil
	}else{
		return ErrCmdDataLengthWrong
	}
}
