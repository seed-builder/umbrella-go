package network

import "encoding/binary"

const(
	CmdConnectReqPktLen uint32 = 12 + 11
	CmdConnectRspPktLen uint32 = 12 + 1

	ConnectWrongSn uint8 = 2
	ConnectSuccess uint8 = 1
	ConnectFail uint8 = 0
)


//CmdConnectReqPkt is the connect request packet
type CmdConnectReqPkt struct {
	EquipmentSn string

	//session info
	SeqId uint32
}

type CmdConnectRspPkt struct {
	Status uint8
	//session info
	SeqId uint32
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdConnectReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdConnectReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CONNECT)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	w.WriteFixedSizeString(p.EquipmentSn, 11)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdConnectReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	sn := r.ReadCString(11)
	p.EquipmentSn = string(sn)
	return r.Error()
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdConnectRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdConnectRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CONNECT_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId
	w.WriteByte(p.Status)

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdConnectRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Status = r.ReadByte()
	return r.Error()
}
