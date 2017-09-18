package network

import (
	"umbrella/utilities"
)

type CmdData struct {
	Length uint8
	SeqId uint8
	CmdId uint8
	Channel uint8
	Body []byte
	Sign uint8
	//数据状态
	DataStatus uint8
	Err error
}

// Pack packs the CmdActiveTestReqPkt to bytes stream for client side.
func (cmd CmdData) ToBytes(body... byte) ([]byte, error) {
	pktLen := 5
	bodyLen := 0
	if body != nil {
		bodyLen = len(body)
		pktLen += bodyLen
		cmd.Body = body
	}
	var w = newPacketWriter(uint32(pktLen))

	// Pack header
	cmd.Length = uint8(pktLen)
	w.WriteByte(cmd.Length)
	w.WriteByte(cmd.SeqId)
	w.WriteByte(cmd.CmdId)
	w.WriteByte(cmd.Channel)
	if bodyLen > 0 {
		w.WriteBytes(cmd.Body)
	}
	buf, _ := w.Bytes()
	var sum byte
	for _, d := range buf {
		sum += d
	}
	cmd.Sign = sum
	w.WriteByte(cmd.Sign)
	return w.Bytes()
}

//解析接收到的命令数据（已除去头尾标识）
func (cmd CmdData) ParseCmdData(buf []byte) (Packer, error) {
	cmd.Length = uint8(buf[0])
	length := len(buf)
	if int(cmd.Length) != length {
		cmd.DataStatus = utilities.RspStatusDataErr
		cmd.Err = ErrCmdIllegal
		return nil, cmd.Err
	}
	data := buf[:length - 1]
	var sum uint8
	for _, d := range data {
		sum += uint8(d)
	}
	crc := uint8(buf[length - 1])
	if crc != sum {
		cmd.DataStatus = utilities.RspStatusDataErr
		cmd.Err = ErrCmdIllegal
		return nil, cmd.Err
	}
	cmd.SeqId = uint8(buf[1])
	// Command_Id
	cmd.CmdId = uint8(buf[2])
	cmd.Channel = uint8(buf[3])
	cmd.Sign = crc
	cmd.Body = buf[4:length - 1]
	//初始化状态为：成功
	cmd.DataStatus = utilities.RspStatusSuccess

	utilities.SysLog.Infof("命令【%x】总长度【%d】,编号【%x】【%s】, 序列号【%d】 ", buf, length, cmd.CmdId, CmdDesc(cmd.CmdId), cmd.SeqId)
	// The left packet data (start from seqId in header).

	var p Packer
	canUnpack := true
	switch cmd.CmdId {
	case CMD_ACTIVE_TEST:
		p = &CmdActiveTestReqPkt{
			CmdData: cmd,
		}
	case CMD_ACTIVE_TEST_RESP:
		p = &CmdActiveTestRspPkt{
			CmdData: cmd,
		}
	case CMD_CONNECT:
		p = &CmdConnectReqPkt{
			CmdData: cmd,
		}
	case CMD_CONNECT_RESP:
		p = &CmdConnectRspPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_INSPECT:
		p = &CmdChannelInspectReqPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_INSPECT_RESP:
		p = &CmdChannelInspectRspPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_TAKE_UMBRELLA:
		p = &CmdTakeUmbrellaReqPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_TAKE_UMBRELLA_RESP:
		p = &CmdTakeUmbrellaRspPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_UMBRELLA_OUT:
		p = &CmdUmbrellaOutReqPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_UMBRELLA_OUT_RESP:
		p = &CmdUmbrellaOutRspPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_UMBRELLA_IN:
		p = &CmdUmbrellaInReqPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_UMBRELLA_IN_RESP:
		p = &CmdUmbrellaInRspPkt{
			CmdData: cmd,
		}
	case CMD_UMBRELLA_INSPECT:
		p = &CmdUmbrellaInspectReqPkt{
			CmdData: cmd,
		}
	case CMD_UMBRELLA_INSPECT_RESP:
		p = &CmdUmbrellaInspectRspPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_RESCUE:
		p = &CmdChannelRescueReqPkt{
			CmdData: cmd,
		}
	case CMD_CHANNEL_RESCUE_RESP:
		p = &CmdChannelRescueRspPkt{
			CmdData: cmd,
		}

	default:
		p = nil
		canUnpack = false
		cmd.DataStatus = utilities.RspStatusCmdIllegal
		cmd.Err = ErrCmdIllegal
		return nil, cmd.Err
	}
	if canUnpack && (length-1) >= 4 {
		err := p.Unpack(buf[4:length-1])
		if err != nil {
			utilities.SysLog.Warningf(" 解析命令详情数据错误： %v ", err)
			cmd.DataStatus = utilities.RspStatusDataErr
			cmd.Err = err
			//return nil, cmd.Err
		}
	}
	return p, nil
}
