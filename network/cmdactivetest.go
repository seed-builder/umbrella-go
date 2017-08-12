package network

import "log"

// Packet length const for cmd active test request and response packets.
const (
	CmdActiveTestReqPktLen uint32 = 4     //12d, 0xc
	CmdActiveTestRspPktLen uint32 = 4 //13d, 0xd
)


type CmdActiveTestReqPkt struct {
	// session info
	SeqId uint8
}

type CmdActiveTestRspPkt struct {
	//Reserved uint8
	// session info
	SeqId uint8
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdActiveTestReqPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdActiveTestReqPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteByte(4)
	w.WriteByte(seqId)
	w.WriteByte(byte(CMD_ACTIVE_TEST))

	p.SeqId = seqId
	log.Println("cmd active test pack seqid: ", seqId)
	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdActiveTestReqPkt) Unpack(data []byte) error {
	//var r = newPacketReader(data)
	////
	////// Sequence Id
	//p.SeqId = r.ReadByte()
	//return r.Error()
	return nil
}

// Pack packs the CmdActiveTestRspPkt to bytes stream for client side.
func (p *CmdActiveTestRspPkt) Pack(seqId uint8) ([]byte, error) {
	var pktLen = CmdActiveTestRspPktLen

	var w = newPacketWriter(pktLen)

	w.WriteByte(4)
	w.WriteByte(seqId)
	p.SeqId = seqId

	w.WriteByte(byte(CMD_ACTIVE_TEST_RESP))

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestRspPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestRspPkt struct.
func (p *CmdActiveTestRspPkt) Unpack(data []byte) error {
	//var r = newPacketReader(data)
	////
	////// Sequence Id
	//////r.ReadInt(binary.BigEndian, &p.SeqId)
	//p.SeqId = r.ReadByte()
	//return r.Error()
	return nil
}
