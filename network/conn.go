package network

import (
	"net"
	"errors"
	"time"
	"umbrella/models"
	"log"
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
	ErrConnIsClosed = errors.New("connection is closed")
	ErrCmdIllegal = errors.New("illegal cmd")
)

var noDeadline = time.Time{}
type UmbrellaInStatus uint8

// Conn States
const (
	CONN_CLOSED State = iota
	CONN_CONNECTED
	CONN_AUTHOK
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
	}
	tc := c.rw.(*net.TCPConn) // Always tcpconn
	tc.SetKeepAlive(true) //Keepalive as default
	return c
}

// Serve a new connection.
func (c *Conn) Serve() {
	defer func() {
		if err := recover(); err != nil {
			c.server.ErrorLog.Printf("panic serving %v: %v\n", c.rw.RemoteAddr(), err)
		}
	}()

	defer c.Close()
	c.server.ErrorLog.Printf("conn serving %v \n", c.rw.RemoteAddr())

	c.server.ErrorLog.Printf("waiting for receiving data: %v  \n", c.rw.RemoteAddr())
	for {
		select {
		case <-c.exceed:
			return // close the connection.
		default:
		}

		rs, err := c.readPacket()
		if err != nil {
			c.server.ErrorLog.Printf("receiving data: %v \n", err)
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue
			}
			if err == ErrCmdIllegal {
				continue
			}
			break
			//continue
		}

		c.server.ErrorLog.Printf("readPacket response has : %d \n", len(rs))
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
		c.server.ErrorLog.Printf("receive a cmd connect request from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

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
		c.server.ErrorLog.Printf("receive a cmd open channel request from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

	case *CmdUmbrellaOutRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}
		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmd open channel response from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

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
		c.server.ErrorLog.Printf("receive a cmd umbrella in request from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

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
		c.server.ErrorLog.Printf("receive a cmd active request from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

	case *CmdActiveTestRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmd active response from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

	case *CmdIllegalRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: p,
		}
		c.server.ErrorLog.Printf("receive a illegal cmd request from %v[%d]\n",
			c.rw.RemoteAddr(), p.SeqId)

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
	return c.SendPkt(r.Packer, r.SeqId)
}

func (c *Conn) startActiveTest(){
	exceed := make(chan struct{})
	c.exceed = exceed
	c.server.ErrorLog.Printf("start active test serving %v \n", c.rw.RemoteAddr())
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
					c.server.ErrorLog.Printf("no client active test response returned from %v for %d times!",
						c.rw.RemoteAddr(), c.n)
					exceed <- struct{}{}
					break
				}
				// send a active test packet to peer, increase the active test counter
				p := &CmdActiveTestReqPkt{}
				err := c.SendPkt(p, <- c.SeqId)
				c.server.ErrorLog.Printf("sending active test to %v \n", c.rw.RemoteAddr())
				if err != nil {
					c.server.ErrorLog.Printf("send cmd active test request to %v error: %v", c.rw.RemoteAddr(), err)
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
		c.server.ErrorLog.Printf("close connection with %v!\n", c.rw.RemoteAddr())
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

	log.Printf("conn send data, len = %d, data = %x \n", len(buf), buf )

	_, err = c.rw.Write(buf) //block write

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
		t := time.Now().Add(timeout)
		c.rw.SetReadDeadline(t)
		defer c.rw.SetReadDeadline(noDeadline)
	}

	leftData := make([]byte, defaultReadBufferSize)
	length, err := c.rw.Read(leftData) //io.ReadFull(c.Conn, leftData)
	if err != nil {
		return nil, err
	}
	log.Printf(" RecvAndUnpackPkt receive client data len = %d, data=%x . \n", length, leftData[:length])

	cmds := c.ParsePkt(length, leftData)
	num := len(cmds)
	log.Printf("ParsePkt cmd len = %d \n.", num)
	var packers []Packer
	if num > 0 {
		for _, cmd := range cmds {
			log.Printf("ParsePkt cmd data = %x .\n", cmd)
			p, err := c.CheckAndUnpackPkt(cmd)
			if err == nil {
				packers = append(packers, p)
			}else{
				log.Printf("CheckAndUnpackPkt err = %v .\n", err)
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

	log.Printf(" CheckAndUnpackPkt receive data total len: %d,  command id : %x \n", length, commandId)
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
			log.Printf(" CheckAndUnpackPkt Unpack data err: %v \n", err)
			return nil, err
		}
	}
	return p, nil
}

