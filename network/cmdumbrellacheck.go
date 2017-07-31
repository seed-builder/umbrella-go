package network


const (
	CmdUmbrellaCheckReqPktLen uint32 = 4 + 8
	CmdUmbrellaCheckRspPktLen uint32 = 4 + 1
)

type CmdUmbrellaCheckReqPkt struct{
	//8字节
	UmbrellaSn string

	SeqId uint8
}

type CmdUmbrellaCheckRspPkt struct{
	Status uint8
	SeqId uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaCheckReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaCheckReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x0C)
	w.WriteByte(byte(CMD_UMBRELLA_CHECK))
	w.WriteByte(seqId)
	p.SeqId = seqId
	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteFixedSizeString(p.UmbrellaSn, UmbrellaSnLen)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaCheckReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	sn := r.ReadCString(UmbrellaSnLen)
	p.UmbrellaSn = string(sn)

	return r.Error()
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaCheckRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaCheckRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x05)
	w.WriteByte(byte(CMD_UMBRELLA_CHECK_RESP))
	w.WriteByte(seqId)
	p.SeqId = seqId
	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaCheckRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	p.Status = r.ReadByte()

	return r.Error()
}
