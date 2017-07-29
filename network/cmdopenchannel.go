package network

const (
	CmdOpenChannelReqPktLen uint32 = 4 + 1
	CmdOpenChannelRspPktLen uint32 = 4 + 1
)

type CmdOpenChannelReqPkt struct{
	//EquipmentSn string
	ChannelNum uint8
	//session info
	SeqId uint8
}

type CmdOpenChannelRspPkt struct{
	Status uint8
	//session info
	SeqId uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdOpenChannelReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdOpenChannelReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x05)
	w.WriteByte(byte(CMD_OPEN_CHANNEL))
	w.WriteByte(seqId)
	p.SeqId = seqId
	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteByte(p.ChannelNum)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdOpenChannelReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	p.ChannelNum = r.ReadByte()

	return r.Error()
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdOpenChannelRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdOpenChannelRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x05)
	w.WriteByte(byte(CMD_OPEN_CHANNEL_RESP))
	w.WriteByte(seqId)

	p.SeqId = seqId
	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdOpenChannelRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	p.Status = r.ReadByte()

	return r.Error()
}