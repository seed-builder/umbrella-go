package network

import "encoding/binary"

const (
	CmdOpenChannelReqPktLen uint32 = 12 + 11 + 1
	CmdOpenChannelRspPktLen uint32 = 12 + 11 + 1 + 1
)

type CmdOpenChannelReqPkt struct{
	//EquipmentSn string
	ChannelNum uint8
	//session info
	SeqId uint32
}

type CmdOpenChannelRspPkt struct{
	Status uint8
	//session info
	SeqId uint32
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdOpenChannelReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdOpenChannelReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CONNECT)
	w.WriteInt(binary.BigEndian, seqId)
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
	r.ReadInt(binary.BigEndian, &p.SeqId)
	//sn := r.ReadCString(11)
	//p.EquipmentSn = string(sn)
	p.ChannelNum = r.ReadByte()

	return r.Error()
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdOpenChannelRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdOpenChannelRspPktLen

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
func (p *CmdOpenChannelRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Status = r.ReadByte()

	return r.Error()
}