package network

const (
	CmdUmbrellaOutReqPktLen uint32 = 4 + 1
	CmdUmbrellaOutRspPktLen uint32 = 4 + 1 + 8
	UmbrellaSnLen int = 8
)

type CmdUmbrellaOutReqPkt struct{
	//session info
	SeqId uint8
	//EquipmentSn string
	ChannelNum uint8
}

type CmdUmbrellaOutRspPkt struct{
	//session info
	SeqId uint8
	Status uint8
	//len 8
	UmbrellaSn string
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaOutReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x05)
	w.WriteByte(byte(CMD_UMBRELLA_OUT))
	w.WriteByte(seqId)
	p.SeqId = seqId
	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteByte(p.ChannelNum)

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

	return r.Error()
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaOutRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(0x0D)
	w.WriteByte(byte(CMD_UMBRELLA_OUT_RESP))
	w.WriteByte(seqId)

	p.SeqId = seqId
	w.WriteByte(p.Status)
	w.WriteFixedSizeString(p.UmbrellaSn, UmbrellaSnLen)

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
	sn := r.ReadCString(UmbrellaSnLen)
	p.UmbrellaSn = string(sn)
	return r.Error()
}