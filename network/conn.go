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
)

var noDeadline = time.Time{}
type UmbrellaInStatus uint8

// Conn States
const (
	CONN_CLOSED State = iota
	CONN_CONNECTED
	CONN_AUTHOK
)

const (
	CMD_REQUEST_MIN, CMD_RESPONSE_MIN CommandId = iota, 0x80 + iota
	CMD_ACTIVE_TEST, CMD_ACTIVE_TEST_RESP
	CMD_CONNECT, CMD_CONNECT_RESP
	CMD_UMBRELLA_OUT, CMD_UMBRELLA_OUT_RESP
	CMD_UMBRELLA_IN, CMD_UMBRELLA_IN_RESP
	CMD_ILLEGAL, CMD_IILEGAL_RESP
	CMD_REQUEST_MAX, CMD_RESPONSE_MAX
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
		var i uint8
		for {
			select {
			case out <- i:
				if i >= 255 {
					i = 0
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
			utilities.SysLog.Panicf("客户端会话严重错误 %v: %v", c.rw.RemoteAddr(), err)
		}
	}()
	defer c.Close()

	utilities.SysLog.Infof("开启客户端【%v】会话, 等待接收命令 ", c.rw.RemoteAddr())
	for {
		select {
		case <-c.exceed:
			utilities.SysLog.Warningf("关闭客户端【%v】会话 ", c.rw.RemoteAddr())
			return // close the connection.
		default:
		}
		//读取命令
		rs, err := c.readPacket()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				utilities.SysLog.Warningf("读取命令超时：%v ", err)
				continue
			}
			if err == ErrCmdIllegal {
				utilities.SysLog.Warning("非法命令")
				continue
			}
			utilities.SysLog.Errorf("读取命令错误：%v ", err)
			break
			//continue
		}

		utilities.SysLog.Infof("客户端【%v】,有【%d】条命令待处理", c.rw.RemoteAddr(), len(rs))
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
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *CmdUmbrellaOutReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaOutRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}

	case *CmdUmbrellaOutRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
		}

	case *CmdUmbrellaInReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaInRspPkt{
				SeqId: p.SeqId,
				ChannelNum: p.ChannelNum,
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
			Packer: &CmdActiveTestRspPkt{},
			SeqId: p.SeqId,
		}
	case *CmdActiveTestRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
		}
	case *CmdIllegalRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: p,
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
	utilities.SysLog.Infof("预备向客户端【%v】发送响应命令", c.rw.RemoteAddr())
	return c.SendPkt(r.Packer, r.SeqId)
}

func (c *Conn) startActiveTest(){
	exceed := make(chan struct{})
	c.exceed = exceed
	utilities.SysLog.Infof("预备向客户端【%v】发送维持包数据 ", c.rw.RemoteAddr())
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
					utilities.SysLog.Infof("没接收到客户端【%v】的维持包反馈【%d】次!",
						c.rw.RemoteAddr(), c.n)
					exceed <- struct{}{}
					break
				}
				// send a active test packet to peer, increase the active test counter
				p := &CmdActiveTestReqPkt{}
				err := c.SendPkt(p, <- c.SeqId)
				utilities.SysLog.Infof("向客户端【%v】发送维持包数据 ", c.rw.RemoteAddr())
				if err != nil {
					utilities.SysLog.Infof("向客户端【%v】发送维持包数据错误【$s】", c.rw.RemoteAddr(), err)
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
		utilities.SysLog.Warningf("关闭客户端【%v】连接!", c.rw.RemoteAddr())
		if c.Equipment != nil {
			c.Equipment.Offline()
		}
		close(c.done)  // let the SeqId goroutine exit.
		c.rw.Close() // close the underlying net.Conn
		c.State = CONN_CLOSED
	}
}

func (c *Conn) SetState(state State) {
	c.State = state
}

func (c *Conn) SetEquipment(equipment *models.Equipment){
	c.Equipment = equipment
}

// SendPkt pack the CMD packet structure and send it to the other peer.
func (c *Conn) SendPkt(packet Packer, seqId uint8) error {
	if c.State == CONN_CLOSED {
		return ErrConnIsClosed
	}
	var buf []byte
	buf = append(buf, 0xAA)
	data, err := packet.Pack(seqId)
	if err != nil {
		return err
	}
	//add content
	var sum byte
	for _, d := range data {
		sum += d
		buf = append(buf, d)
	}
	//add crc
	buf = append(buf, sum)
	buf = append(buf, 0x55)

	_, err = c.rw.Write(buf) //block write
	utilities.SysLog.Infof("发送命令【%x】序列号【%d】", buf[:], seqId)
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
	utilities.SysLog.Infof("读取到的数据【%x】长度【%d】", leftData[:length], length)

	cmds := c.ParsePkt(length, leftData)
	num := len(cmds)
	utilities.SysLog.Infof("解析出【%d】条命令 .", num)
	var packers []Packer
	if  num > 0 {
		for _, cmd := range cmds {
			utilities.SysLog.Infof("解析命令详情数据【%x】.", cmd)
			p, err := c.CheckAndUnpackPkt(cmd)
			if err == nil {
				packers = append(packers, p)
			}else{
				utilities.SysLog.Infof("解析命令详情数据【%x】错误： %v .", cmd, err)
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

func (c *Conn) CheckAndUnpackPkt(leftData []byte) (Packer, error)  {
	length := len(leftData)
	data := leftData[:length - 1]
	var sum uint8
	for _, d := range data {
		sum += uint8(d)
	}
	crc := uint8(leftData[length - 1])
	if crc != sum {
		return nil, ErrCmdIllegal
	}
	//totalLen := uint8(leftData[1])
	//seq id
	seqId := uint8(leftData[1])
	// Command_Id
	commandId := CommandId(leftData[2])

	utilities.SysLog.Infof("命令【%x】总长度【%d】,编号【%x】【%s】, 序列号【%d】 ", leftData, length, commandId, CmdDesc(commandId), seqId)
	// The left packet data (start from seqId in header).

	var p Packer
	canUnpack := true
	switch commandId {
	case CMD_CONNECT:
		p = &CmdConnectReqPkt{
			SeqId: seqId,
		}
	case CMD_CONNECT_RESP:
		p = &CmdConnectRspPkt{
			SeqId: seqId,
		}
	case CMD_UMBRELLA_OUT:
		p = &CmdUmbrellaOutReqPkt{
			SeqId: seqId,
		}
	case CMD_UMBRELLA_OUT_RESP:
		p = &CmdUmbrellaOutRspPkt{
			SeqId: seqId,
		}
	case CMD_UMBRELLA_IN:
		p = &CmdUmbrellaInReqPkt{
			SeqId: seqId,
		}
	case CMD_UMBRELLA_IN_RESP:
		p = &CmdUmbrellaInRspPkt{
			SeqId: seqId,
		}
	case CMD_ACTIVE_TEST:
		p = &CmdActiveTestReqPkt{
			SeqId: seqId,
		}
	case CMD_ACTIVE_TEST_RESP:
		p = &CmdActiveTestRspPkt{
			SeqId: seqId,
		}
	default:
		p = &CmdIllegalRspPkt{
			SeqId: seqId,
		}
		canUnpack = false
		//return nil, ErrCommandIdNotSupported
	}
	if canUnpack && (length-1) > 3 {
		err := p.Unpack(leftData[3:length-1])
		if err != nil {
			utilities.SysLog.Warningf(" 解析命令详情数据错误： %v ", err)
			return nil, err
		}
	}
	return p, nil
}

func CmdDesc(cmdId CommandId) string {
	var desc string
	switch cmdId {
	case CMD_CONNECT:
		desc = "登陆请求"
	case CMD_CONNECT_RESP:
		desc = "登陆响应"
	case CMD_UMBRELLA_OUT:
		desc = "出伞请求"
	case CMD_UMBRELLA_OUT_RESP:
		desc = "出伞响应"
	case CMD_UMBRELLA_IN:
		desc = "进伞请求"
	case CMD_UMBRELLA_IN_RESP:
		desc = "进伞响应"
	case CMD_ACTIVE_TEST:
		desc = "维持包请求"
	case CMD_ACTIVE_TEST_RESP:
		desc = "维持包响应"
	}
	return desc
}

