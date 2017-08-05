package network


const(
	CmdIllegalRspPktLen uint32 = 4 + 1
)

type CmdIllegalRspPkt struct {
	//session info
	SeqId uint8
	Status uint8
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdIllegalRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdIllegalRspPktLen

	var w = newPacketWriter(pktLen)

	w.WriteByte(0x05)
	w.WriteByte(seqId)
	p.SeqId = seqId
	w.WriteByte(byte(CMD_IILEGAL_RESP))

	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdIllegalRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	//p.SeqId = r.ReadByte()
	p.Status = r.ReadByte()
	return r.Error()
}
