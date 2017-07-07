package network

import (
	"encoding/binary"
)

const (
	CmdUmbrellaInReqPktLen uint32 = 12 + 7 + 1
	CmdUmbrellaInRspPktLen uint32 = 12 + 1

	UmbrellaSnLen int = 7
)

type CmdUmbrellaInReqPkt struct{

	ChannelNum uint8
	//7字节
	UmbrellaSn string

	SeqId uint32
}

type CmdUmbrellaInRspPkt struct{
	Status uint8
	SeqId uint32
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaInReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdUmbrellaInReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_UMBRELLA_IN)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	//w.WriteFixedSizeString(p.EquipmentSn, 11)
	w.WriteByte(p.ChannelNum)
	w.WriteFixedSizeString(p.UmbrellaSn, UmbrellaSnLen)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaInReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	p.ChannelNum = r.ReadByte()
	usn := r.ReadCString(UmbrellaSnLen)
	p.UmbrellaSn = string(usn)

	return r.Error()
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdUmbrellaInRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdUmbrellaInRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_UMBRELLA_IN_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdUmbrellaInRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Status = r.ReadByte()

	return r.Error()
}
