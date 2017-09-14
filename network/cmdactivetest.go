package network

import "umbrella/utilities"

type CmdActiveTestReqPkt struct {
	CmdData
}

type CmdActiveTestRspPkt struct {
	CmdData
	Status uint8
}


// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (p *CmdActiveTestReqPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_ACTIVE_TEST
	return p.ToBytes()
}

// Unpack unpack the binary byte stream to a CmdActiveTestReqPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestReqPkt struct.
func (p *CmdActiveTestReqPkt) Unpack(data []byte) error {
	return nil
}

// Pack packs the CmdActiveTestRspPkt to bytes stream for client side.
func (p *CmdActiveTestRspPkt) Pack(seqId uint8) ([]byte, error) {
	p.SeqId = seqId
	p.CmdId = CMD_ACTIVE_TEST_RESP
	p.Status = utilities.RspStatusSuccess
	return p.ToBytes(p.Status)
}

// Unpack unpack the binary byte stream to a CmdActiveTestRspPkt variable.
// After unpack, you will get all value of fields in
// CmdActiveTestRspPkt struct.
func (p *CmdActiveTestRspPkt) Unpack(data []byte) error {
	p.Status = data[0]
	return nil
}
