package network

import "encoding/binary"

const (
	CmdTerminateReqPktLen uint32 = 12 //12d, 0xc
	CmdTerminateRspPktLen uint32 = 12 //12d, 0xc
)

type CmdTerminateReqPkt struct{
	EquipmentSn string
	SeqId uint32
}

type CmdTerminateRspPkt struct{
	EquipmentSn string
	SeqId uint32
}

// Pack packs the CmppTerminateReqPkt to bytes stream for client side.
func (p *CmdTerminateReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdTerminateReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_TERMINATE)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	w.WriteFixedSizeString(p.EquipmentSn, 11)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateReqPkt variable.
// After unpack, you will get all value of fields in
// CmppTerminateReqPkt struct.
func (p *CmdTerminateReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	sn := r.ReadCString(11)
	p.EquipmentSn = string(sn)
	return r.Error()
}

// Pack packs the CmppTerminateRspPkt to bytes stream for client side.
func (p *CmdTerminateRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdTerminateRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_TERMINATE_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	w.WriteFixedSizeString(p.EquipmentSn, 11)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmppTerminateRspPkt variable.
// After unpack, you will get all value of fields in
// CmppTerminateRspPkt struct.
func (p *CmdTerminateRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	sn := r.ReadCString(11)
	p.EquipmentSn = string(sn)
	return r.Error()
}
