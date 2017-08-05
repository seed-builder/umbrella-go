package network


const(
	CmdConnectReqPktLen uint32 = 4 + 11
	CmdConnectRspPktLen uint32 = 4 + 1

	ConnectWrongSn uint8 = 2
	ConnectSuccess uint8 = 1
	ConnectFail uint8 = 0

	EquipmentSnLen int = 11

)


//CmdConnectReqPkt is the connect request packet
type CmdConnectReqPkt struct {
	//session info
	SeqId uint8
	EquipmentSn string
}

type CmdConnectRspPkt struct {
	//session info
	SeqId uint8
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdConnectReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdConnectReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x0F)
	w.WriteByte(seqId)
	p.SeqId = seqId
	w.WriteByte(byte(CMD_CONNECT))

	w.WriteFixedSizeString(p.EquipmentSn, EquipmentSnLen)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdConnectReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)
	// Sequence Id
	//p.SeqId = r.ReadByte()
	sn := r.ReadCString(EquipmentSnLen)
	p.EquipmentSn = string(sn)
	return r.Error()
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdConnectRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdConnectRspPktLen

	var w = newPacketWriter(pktLen)

	w.WriteByte(0x05)
	w.WriteByte(seqId)
	p.SeqId = seqId
	w.WriteByte(byte(CMD_CONNECT_RESP))

	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdConnectRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	//p.SeqId = r.ReadByte()
	p.Status = r.ReadByte()
	return r.Error()
}
