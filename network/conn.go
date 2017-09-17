package network

import (
	"net"
	"errors"
	"time"
	"umbrella/models"
	"fmt"
	"sync/atomic"
	"umbrella/utilities"
)

type State uint8
type Type int8

const (
	V10 Type = 0x10
	CMDHEAD byte = 0xAA
	CMDFOOT byte = 0x55
	defaultReadBufferSize = 64
)

// Errors for conn operations
var (
	ErrConnIsClosed = errors.New("连接已关闭")
	ErrCmdIllegal = errors.New("非法命令")
	ErrCmdDataLengthWrong = errors.New("数据长度错误")
)

var noDeadline = time.Time{}
type UmbrellaInStatus uint8

// Conn States
const (
	CONN_CLOSED State = iota
	CONN_CONNECTED
	CONN_AUTHOK
)

//命令
const (
	CMD_ACTIVE_TEST uint8 = 0x01
	CMD_ACTIVE_TEST_RESP uint8 = 0x81
	CMD_CONNECT uint8 = 0x02
    CMD_CONNECT_RESP uint8 = 0x82
	CMD_CHANNEL_INSPECT uint8 = 0x41
	CMD_CHANNEL_INSPECT_RESP uint8 = 0xC1
	CMD_CHANNEL_TAKE_UMBRELLA uint8 = 0x42
	CMD_CHANNEL_TAKE_UMBRELLA_RESP uint8 = 0xC2
	CMD_CHANNEL_UMBRELLA_OUT uint8 = 0x43
	CMD_CHANNEL_UMBRELLA_OUT_RESP uint8 = 0xC3
	CMD_CHANNEL_UMBRELLA_IN uint8 = 0x44
	CMD_CHANNEL_UMBRELLA_IN_RESP uint8 = 0xC4
	CMD_UMBRELLA_INSPECT uint8 = 0x45
	CMD_UMBRELLA_INSPECT_RESP uint8 = 0xC5
)

type Conn struct {
	server *TcpServer // the Server on which the connection arrived

	// for active test
	t       time.Duration // interval between two active tests
	n       int32         // continuous send times when no response back
	done    chan struct{}
	exceed  chan struct{}
	counter int32

    rw	net.Conn
	State State
	Typ   Type

	// for SeqId generator goroutine
	SeqId <-chan uint8
	//done  chan<- struct{}

	Equipment *models.Equipment
	Ip string
}

func newSeqIdGenerator() (<-chan uint8, chan struct{}) {
	out := make(chan uint8)
	done := make(chan struct{})

	go func() {
		var i uint8 = 1
		for {
			select {
			case out <- i:
				if i >= 255 {
					i = 1
				}else{
					i ++
				}
			case <-done:
				close(out)
				return
			}
		}
	}()
	return out, done
}

// New returns an abstract structure for successfully
// established underlying net.Conn.
func NewConn(svr *TcpServer, conn net.Conn, typ Type) *Conn {
	seqId, done := newSeqIdGenerator()
	c := &Conn{
		server: svr,
		rw:  conn,
		Typ:   typ,
		SeqId: seqId,
		done:  done,
		Ip: conn.RemoteAddr().String(),
	}
	tc := c.rw.(*net.TCPConn) // Always tcpconn
	tc.SetKeepAlive(true) //Keepalive as default
	return c
}

// Serve a new connection.
func (c *Conn) Serve() {
	defer func() {
		if err := recover(); err != nil {
			c.Panicf("客户端会话严重错误 %v: %v", c.rw.RemoteAddr(), err)
		}
	}()
	defer c.Close()

	c.Noticef("开启客户端【%v】会话, 等待接收命令 ", c.rw.RemoteAddr())
	for {
		select {
		case <-c.exceed:
			c.Warningf("关闭客户端【%v】会话 ", c.rw.RemoteAddr())
			return // close the connection.
		default:
		}
		//读取命令
		rs, err := c.readPacket()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				c.Debugf("读取命令超时：%v ", err)
				continue
			}
			c.Errorf("读取命令错误：%v ", err)
			break
			//continue
		}

		c.Infof("客户端【%v】,有【%d】条命令待处理", c.rw.RemoteAddr(), len(rs))
		for _, r := range rs {
			_, err = c.server.Handler.ServeHandle(r, r.Packet, c.server.ErrorLog)
			if err1 := c.finishPacket(r); err1 != nil {
				break
			}

			if err != nil {
				break
			}
		}
	}
}

func (c *Conn) readPacket() ([]*Response, error) {
	readTimeout := time.Second * 5
	packers, err := c.RecvAndUnpackPkt(readTimeout)
	if err != nil {
		return nil, err
	}
	//typ := c.server.Typ
	var rsps []*Response
	for _, p := range packers {
		r, e := c.parseResponse(p)
		if e == nil {
			rsps = append(rsps, r)
		}
	}
	return rsps, nil
}

func (c *Conn) parseResponse(i Packer) (*Response, error)  {
	var pkt *Packet
	var rsp *Response
	switch p := i.(type) {
	case *CmdConnectReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdConnectRspPkt{
				CmdData: CmdData{
					SeqId: p.SeqId,
				},
			},
			SeqId: p.SeqId,
		}
	case *CmdUmbrellaInReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaInRspPkt{
				CmdData: CmdData{
					SeqId: p.SeqId,
					Channel: p.Channel,
				},
				UmbrellaSn: p.UmbrellaSn,
				Status: utilities.RspStatusSuccess,
			},
			SeqId: p.SeqId,
		}
	case *CmdActiveTestReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdActiveTestRspPkt{
				CmdData: CmdData{
					SeqId: p.SeqId,
					Channel: p.Channel,
				},
				Status: utilities.RspStatusSuccess,
			},
			SeqId: p.SeqId,
		}
	case *CmdUmbrellaInspectReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaInspectRspPkt{
				CmdData: CmdData{
					SeqId: p.SeqId,
					Channel: p.Channel,
				},
				UmbrellaSn: p.UmbrellaSn,
				Status: utilities.RspStatusSuccess,
			},
			SeqId: p.SeqId,
		}
	case *CmdTakeUmbrellaRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaOutReqPkt{
				CmdData: CmdData{
					SeqId: p.SeqId,
					Channel: p.Channel,
				},
				UmbrellaSn: p.UmbrellaSn,
			},
			SeqId: p.SeqId,
		}
	case *CmdChannelInspectRspPkt,*CmdUmbrellaOutRspPkt,*CmdActiveTestRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
		}
	default:
		return nil, NewOpError(ErrUnsupportedPkt,
			fmt.Sprintf("readPacket: receive unsupported packet type: %#v", p))
	}
	return rsp, nil
}

func (c *Conn) finishPacket(r *Response) error {
	if _, ok := r.Packet.Packer.(*CmdActiveTestRspPkt); ok {
		atomic.AddInt32(&c.counter, -1)
		return nil
	}

	if r.Packer == nil {
		// For response packet received, it need not
		// to send anything back.
		return nil
	}
	if rsp, ok := r.Packer.(*CmdConnectRspPkt); ok && rsp.Status == utilities.RspStatusSuccess{
		// start a goroutine for sending active test.
		c.startActiveTest()
	}
	c.Infof("预备向客户端【%v】发送响应命令", c.rw.RemoteAddr())
	return c.SendPkt(r.Packer, r.SeqId)
}

func (c *Conn) startActiveTest(){
	exceed := make(chan struct{})
	c.exceed = exceed
	c.Infof("预备向客户端【%v】发送维持包数据 ", c.rw.RemoteAddr())
	go func() {
		t := time.NewTicker(c.t)
		defer t.Stop()
		for {
			select {
			case <- c.done:
				// once conn close, the goroutine should exit
				return
			case <- t.C:
				// check whether c.counter exceeds
				if atomic.LoadInt32(&c.counter) >= c.n {
					c.Infof("没接收到客户端【%v】的维持包反馈【%d】次!",
						c.rw.RemoteAddr(), c.n)
					exceed <- struct{}{}
					break
				}
				// send a active test packet to peer, increase the active test counter
				p := &CmdActiveTestReqPkt{}
				err := c.SendPkt(p, 0)
				c.Infof("向客户端【%v】发送维持包数据 ", c.rw.RemoteAddr())
				if err != nil {
					c.Infof("向客户端【%v】发送维持包数据错误【$s】", c.rw.RemoteAddr(), err)
				} else {
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}

func (c *Conn) Close() {
	if c != nil {
		if c.State == CONN_CLOSED {
			return
		}
		c.Warningf("关闭客户端【%v】连接!", c.rw.RemoteAddr())
		if c.Equipment != nil {
			c.Equipment.Offline()
		}
		close(c.done)  // let the SeqId goroutine exit.
		c.rw.Close() // close the underlying net.Conn
		close(c.exceed)
		c.State = CONN_CLOSED
		msg := &models.Message{}
		msg.AddEquipmentError(c.Equipment.Sn, c.Equipment.ID, c.Equipment.SiteId, "设备异常下线")
	}
}

func (c *Conn) SetState(state State) {
	c.State = state
}

func (c *Conn) SetEquipment(equipment *models.Equipment){
	c.Equipment = equipment
	go func(){
		time.Sleep(1*time.Second)
		c.ChannelInspect(1)
	}()
}

func (c *Conn) SetChannelStatus(num uint8, status uint8){
	c.Equipment.SetChannelStatus(num, status)
}

func (c *Conn) ChannelInspect(channel uint8){
		req := &CmdChannelInspectReqPkt{
			CmdData: CmdData{ Channel: channel, },
		}
		seqId := <- c.SeqId
		c.Infof("发送通道检测命令设备【%s】通道【%d】序号【%d】", c.Equipment.Sn, channel, seqId)
		c.SendPkt(req, seqId)
}

// SendPkt pack the CMD packet structure and send it to the other peer.
func (c *Conn) SendPkt(packet Packer, seqId uint8) error {
	if c.State == CONN_CLOSED {
		return ErrConnIsClosed
	}
	var buf []byte
	buf = append(buf, CMDHEAD)
	data, err := packet.Pack(seqId)
	if err != nil {
		return err
	}
	buf = append(buf, data...)
	buf = append(buf, CMDFOOT)

	_, err = c.rw.Write(buf) //block write
	c.Noticef("发送命令【%s】长度【%d】序列号【%d】命令ID【%X】通道【%d】数据【%X】--【%X】", CmdDesc(data[2]), data[0], data[1], data[2], data[3], data[4:], buf[:])
	buf = nil
	if err != nil {
		return err
	}
	return nil
}

// RecvAndUnpackPkt receives CMD byte stream, and unpack it to some CMD packet structure.
func (c *Conn) RecvAndUnpackPkt(timeout time.Duration) ([]Packer, error) {
	if c.State == CONN_CLOSED {
		return nil, ErrConnIsClosed
	}

	if timeout != 0 {
		t := time.Now().Local().Add(timeout)
		c.rw.SetReadDeadline(t)
		defer c.rw.SetReadDeadline(noDeadline)
	}

	leftData := make([]byte, defaultReadBufferSize)
	length, err := c.rw.Read(leftData) //io.ReadFull(c.Conn, leftData)
	if err != nil {
		return nil, err
	}
	c.Infof("读取到的数据【%x】长度【%d】", leftData[:length], length)

	cmds := c.ParsePkt(length, leftData)
	num := len(cmds)
	c.Infof("解析出【%d】条命令 .", num)
	var packers []Packer
	if  num > 0 {
		for _, cmd := range cmds {
			//c.Infof("解析命令详情数据【%x】.", cmd)
			c.Noticef("接收命令详情：【%s】长度【%d】序列号【%d】命令ID【%X】通道【%d】数据【%x】", CmdDesc(cmd[2]), cmd[0], cmd[1], cmd[2], cmd[3], cmd[4:])
			cd := &CmdData{}
			p, err := cd.ParseCmdData(cmd)
			if err == nil {
				packers = append(packers, p)
			}else{
				c.Warningf("解析命令详情数据【%x】错误： %v .", cmd, err)
			}
		}
	}
	return packers, nil
}

func (c *Conn) ParsePkt(len int, data []byte) [][]byte {
	var result [][]byte
	head := -1
	foot := -1
	complete := 0
	i := 0
	for i < len {
		d := data[i]
		if head == -1 && d == CMDHEAD {
			head = i
			dataLen := data[i+1]
			i += int(dataLen) + 1
			complete ++
			continue
		}
		if foot == -1 && d == CMDFOOT {
			foot = i
			complete ++
		}
		if complete == 2 && foot > head {
			result = append(result, data[head+1:foot])
			head = -1
			foot = -1
			complete = 0
		}
		i ++
	}
	return result
}

//log
func (c *Conn) Debugf(format string, args... interface{}){
	msg := c.logMsg(format, args...)
	utilities.SysLog.Debug(msg)
}

func (c *Conn) Infof(format string, args... interface{}){
	msg := c.logMsg(format, args...)
	utilities.SysLog.Info(msg)
}

func (c *Conn) Noticef(format string, args... interface{}){
	msg := c.logMsg(format, args...)
	utilities.SysLog.Noticef(msg)
}

func (c *Conn) Errorf(format string, args... interface{}){
	msg := c.logMsg(format, args...)
	utilities.SysLog.Error(msg)
}

func (c *Conn) Warningf(format string, args... interface{}){
	msg := c.logMsg(format, args...)
	utilities.SysLog.Warning(msg)
}

func (c *Conn) Panicf(format string, args... interface{}){
	msg := c.logMsg(format, args...)
	utilities.SysLog.Panicf(msg)
}

func (c *Conn) logMsg(format string, args... interface{}) string{
	msg := ""
	if c.Equipment != nil{
		msg = fmt.Sprintf("设备【%s】", c.Equipment.Sn) + fmt.Sprintf(format, args...)
	}else{
		msg = fmt.Sprintf(format, args...)
	}
	return msg
}

func CmdDesc(cmdId uint8) string {
	var desc string
	switch cmdId {
	case CMD_CONNECT:
		desc = "登陆请求"
	case CMD_CONNECT_RESP:
		desc = "登陆响应"
	case CMD_CHANNEL_TAKE_UMBRELLA:
		desc = "取伞请求"
	case CMD_CHANNEL_TAKE_UMBRELLA_RESP:
		desc = "取伞响应"
	case CMD_CHANNEL_UMBRELLA_OUT:
		desc = "出伞请求"
	case CMD_CHANNEL_UMBRELLA_OUT_RESP:
		desc = "出伞响应"
	case CMD_CHANNEL_UMBRELLA_IN:
		desc = "进伞请求"
	case CMD_CHANNEL_UMBRELLA_IN_RESP:
		desc = "进伞响应"
	case CMD_ACTIVE_TEST:
		desc = "维持包请求"
	case CMD_ACTIVE_TEST_RESP:
		desc = "维持包响应"
	case CMD_CHANNEL_INSPECT:
		desc = "通道检查"
	case CMD_CHANNEL_INSPECT_RESP:
		desc = "通道检查响应"
	case CMD_UMBRELLA_INSPECT:
		desc = "伞SN检查"
	case CMD_UMBRELLA_INSPECT_RESP:
		desc = "伞SN检查响应"
	}
	return desc
}

