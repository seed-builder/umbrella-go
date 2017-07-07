package network

import "encoding/binary"

const(
	CmdCloseChannelPktLen uint32 = 12 + 1
	CmdCloseChannelRspPktlLen uint32 = 12 + 1
)


type CmdCloseChannelReqPkt struct{
	ChannelNum uint8
	// session info
	SeqId uint32
}

type CmdCloseChannelRspPkt struct{
	Status uint8
	// session info
	SeqId uint32
}

// Pack packs the CmdActiveTestRspPkt to bytes stream for client side.
func (p *CmdCloseChannelReqPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CmdCloseChannelPktLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CLOSE_CHANNEL)
	w.WriteInt(binary.BigEndian, seqId)
	w.WriteByte(p.ChannelNum)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestRspPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestRspPkt struct.
func (p *CmdCloseChannelReqPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.ChannelNum = r.ReadByte()
	return r.Error()
}


// Pack packs the CmdActiveTestRspPkt to bytes stream for client side.
func (p *CmdCloseChannelRspPkt) Pack(seqId uint32) ([]byte, error) {
	var pktLen uint32 = CmdCloseChannelRspPktlLen

	var w = newPacketWriter(pktLen)

	// Pack header
	w.WriteInt(binary.BigEndian, pktLen)
	w.WriteInt(binary.BigEndian, CMD_CLOSE_CHANNEL_RESP)
	w.WriteInt(binary.BigEndian, seqId)
	w.WriteByte(p.Status)
	p.SeqId = seqId

	return w.Bytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestRspPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestRspPkt struct.
func (p *CmdCloseChannelRspPkt) Unpack(data []byte) error {
	var r = newPacketReader(data)

	// Sequence Id
	r.ReadInt(binary.BigEndian, &p.SeqId)
	p.Status = r.ReadByte()

	return r.Error()
}
