package network

const (
	CmdUmbrellaOutReqPktLen uint32 = 4 + 7 + 1
	CmdUmbrellaOutRspPktLen uint32 = 4 + 1
)

type CmdUmbrellaOutReqPkt struct{
	ChannelNum uint8
	//7字节
	UmbrellaSn string

	SeqId uint8
}

type CmdUmbrellaOutRspPkt struct{
	Status uint8
	SeqId uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaOutReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x0C)
	w.WriteByte(byte(CMD_UMBRELLA_OUT))
	w.WriteByte(seqId)
	p.SeqId = seqId
	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteByte(p.ChannelNum)
	w.WriteFixedSizeString(p.UmbrellaSn, UmbrellaSnLen)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaOutReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	p.ChannelNum = r.ReadByte()
	usn := r.ReadCString(UmbrellaSnLen)
	p.UmbrellaSn = string(usn)

	return r.Error()
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaOutRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x05)
	w.WriteByte(byte(CMD_UMBRELLA_OUT_RESP))
	w.WriteByte(seqId)

	p.SeqId = seqId
	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaOutRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	p.SeqId = r.ReadByte()
	p.Status = r.ReadByte()

	return r.Error()
}
