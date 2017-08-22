package network

import "log"

const (
	CmdUmbrellaInReqPktLen uint32 = 5+4
	CmdUmbrellaInRspPktLen uint32 = 5+1

	//UmbrellaSnLen int = 7
)

type CmdUmbrellaInReqPkt struct{
	SeqId uint8
	ChannelNum uint8
	//4字节
	UmbrellaSn []byte
}

type CmdUmbrellaInRspPkt struct{
	SeqId uint8
	ChannelNum uint8
	Status uint8
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaInReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaInReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(byte(CmdUmbrellaInReqPktLen))
	w.WriteByte(seqId)
	p.SeqId = seqId
	w.WriteByte(byte(CMD_UMBRELLA_IN))

	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteByte(p.ChannelNum)
	//w.WriteFixedSizeString(p.UmbrellaSn, UmbrellaSnLen)
	//w.WriteInt(binary.LittleEndian, p.UmbrellaSn)
	for _, b := range p.UmbrellaSn {
		w.WriteByte(b)
	}

	return w.Bytes()

}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaInReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	//p.SeqId = r.ReadByte()
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	p.ChannelNum = r.ReadByte()
	//sn := r.ReadCString(UmbrellaSnLen)
	//p.UmbrellaSn = string(sn)
	//r.ReadInt(binary.LittleEndian, &p.UmbrellaSn)
	for i:=0 ; i < 4; i ++{
		p.UmbrellaSn = append(p.UmbrellaSn, r.ReadByte())
	}
	log.Printf(" CmdUmbrellaInReqPkt Unpack ChannelNum: %d, UmbrellaSn: %x  \n", p.ChannelNum , p.UmbrellaSn)
	return r.Error()
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaInRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdUmbrellaInRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(byte(CmdUmbrellaInRspPktLen))
	w.WriteByte(seqId)
	p.SeqId = seqId
	w.WriteByte(byte(CMD_UMBRELLA_IN_RESP))

	w.WriteByte(p.ChannelNum)
	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaInRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	//p.SeqId = r.ReadByte()
	p.ChannelNum = r.ReadByte()
	p.Status = r.ReadByte()

	return r.Error()
}
