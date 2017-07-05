package network

import "encoding/binary"

// Packet length const for cmd active test request and response packets.
const (
	CmdActiveTestReqPktLen uint32 = 12     //12d, 0xc
	CmdActiveTestRspPktLen uint32 = 12 + 1 //13d, 0xd
)


type CmdActiveTestReqPkt struct {
	// session info
	SeqId uint32
}

type CmdActiveTestRspPkt struct {
	Reserved uint8
	// session info
	SeqId uint32
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdActiveTestReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdActiveTestReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_ACTIVE_TEST)
	w.WriteInt(binary.BigEndian, seqId)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdActiveTestReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	return r.Error()
}

// Pack packs the CmdActiveTestRspPkt to bytes stream for client side.
func (p *CmdActiveTestRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen = CmdActiveTestRspPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_ACTIVE_TEST_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	w.WriteByte(p.Reserved)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestRspPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestRspPkt struct.
func (p *CmdActiveTestRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Reserved = r.ReadByte()
	return r.Error()
}
