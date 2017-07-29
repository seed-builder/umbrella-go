package network

const (
	CmdTerminateReqPktLen uint32 = 4 //12d, 0xc
	CmdTerminateRspPktLen uint32 = 4 //12d, 0xc
)

type CmdTerminateReqPkt struct{
	SeqId uint8
}

type CmdTerminateRspPkt struct{
	SeqId uint8
}

// Pack packs the CmppTerminateReqPkt to bytes stream for client side.
func (p *CmdTerminateReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdTerminateReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x04)
	w.WriteByte(byte(CMD_TERMINATE))
	w.WriteByte(seqId)

	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateReqPkt variable.
// After unpack, you will get all value of fields in
// CmppTerminateReqPkt struct.
func (p *CmdTerminateReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	return r.Error()
}

// Pack packs the CmppTerminateRspPkt to bytes stream for client side.
func (p *CmdTerminateRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdTerminateRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x04)
	w.WriteByte(byte(CMD_TERMINATE_RESP))
	w.WriteByte(seqId)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateRspPkt variable.
// After unpack, you will get all value of fields in
// CmppTerminateRspPkt struct.
func (p *CmdTerminateRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	return r.Error()
}
