package network

import "encoding/binary"

const (
	CmdUmbrellaOutReqPktLen uint32 = 12 + 11 + 1 + 11
	CmdUmbrellaOutRspPktLen uint32 = 12 + 1
)

type UmbrellaOutStatus uint8

const (
	//成功记录借伞
	UmbrellaOutStatusSuccess UmbrellaOutStatus = iota + 1
	UmbrellaOutStatusFail
)


type CmdUmbrellaOutReqPkt struct{
	ChannelNum uint8
	UmbrellaSn string

	SeqId uint32
}

type CmdUmbrellaOutRspPkt struct{
	Status uint8
	SeqId uint32
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdUmbrellaOutReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CONNECT)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteByte(p.ChannelNum)
	w.WriteFixedSizeString(p.UmbrellaSn, 11)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaOutReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	p.ChannelNum = r.ReadByte()
	usn := r.ReadCString(11)
	p.UmbrellaSn = string(usn)

	return r.Error()
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaOutRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdUmbrellaOutRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CONNECT)
	w.WriteInt(binary.BigEndian, seqId)
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
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Status = r.ReadByte()

	return r.Error()
}
