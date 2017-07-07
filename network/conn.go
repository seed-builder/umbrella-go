package network

import (
	"net"
	"errors"
	"time"
	"encoding/binary"
	"io"
	"sync"
	"umbrella/models"
	"log"
)

type State uint8
type Type int8


const (
	V10 Type = 0x10
)

// Errors for conn operations
var (
	ErrConnIsClosed = errors.New("connection is closed")
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
	net.Conn
	State State
	Typ   Type

	// for SeqId generator goroutine
	SeqId <-chan uint32
	done  chan<- struct{}

	Equipment *models.Equipment
}

func newSeqIdGenerator() (<-chan uint32, chan<- struct{}) {
	out := make(chan uint32)
	done := make(chan struct{})

	go func() {
		var i uint32
		for {
			select {
			case out <- i:
				i++
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
func NewConn(conn net.Conn, typ Type) *Conn {
	seqId, done := newSeqIdGenerator()
	c := &Conn{
		Conn:  conn,
		Typ:   typ,
		SeqId: seqId,
		done:  done,
	}
	tc := c.Conn.(*net.TCPConn) // Always tcpconn
	tc.SetKeepAlive(true)       //Keepalive as default
	return c
}

func (c *Conn) Close() {
	if c != nil {
		if c.State == CONN_CLOSED {
			return
		}
		c.Equipment.Offline()
		close(c.done)  // let the SeqId goroutine exit.
		c.Conn.Close() // close the underlying net.Conn
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
func (c *Conn) SendPkt(packet Packer, seqId uint32) error {
	if c.State == CONN_CLOSED {
		return ErrConnIsClosed
	}

	data, err := packet.Pack(seqId)
	if err != nil {
		return err
	}

	_, err = c.Conn.Write(data) //block write
	if err != nil {
		return err
	}

	return nil
}

const (
	defaultReadBufferSize = 4096
)

// readBuffer is used to optimize the performance of
// RecvAndUnpackPkt.
type readBuffer struct {
	totalLen  uint32
	commandId CommandId
	leftData  [defaultReadBufferSize]byte
}

var readBufferPool = sync.Pool{
	New: func() interface{} {
		return &readBuffer{}
	},
}

// RecvAndUnpackPkt receives CMD byte stream, and unpack it to some CMD packet structure.
func (c *Conn) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	if c.State == CONN_CLOSED {
		return nil, ErrConnIsClosed
	}

	if timeout != 0 {
		t := time.Now().Add(timeout)
		c.SetReadDeadline(t)
		defer c.SetReadDeadline(noDeadline)
	}

	rb := readBufferPool.Get().(*readBuffer)
	defer readBufferPool.Put(rb)

	// Total_Length in packet
	err := binary.Read(c.Conn, binary.BigEndian, &rb.totalLen)
	if err != nil {
		return nil, err
	}

	if rb.totalLen < CMD_PACKET_MIN || rb.totalLen > CMD_PACKET_MAX {
		return nil, ErrTotalLengthInvalid
	}


	// Command_Id
	err = binary.Read(c.Conn, binary.BigEndian, &rb.commandId)
	if err != nil {
		return nil, err
	}

	if !((rb.commandId > CMD_REQUEST_MIN && rb.commandId < CMD_REQUEST_MAX) ||
			(rb.commandId > CMD_RESPONSE_MIN && rb.commandId < CMD_RESPONSE_MAX)) {
		return nil, ErrCommandIdInvalid
	}

	log.Println("receive data total len: ", rb.totalLen, " command id : ", rb.commandId)

	// The left packet data (start from seqId in header).
	var leftData = rb.leftData[0:(rb.totalLen - 8)]
	_, err = io.ReadFull(c.Conn, leftData)
	if err != nil {
		return nil, err
	}

	var p Packer
	switch rb.commandId {
	case CMD_CONNECT:
		p = &CmdConnectReqPkt{}
	case CMD_CONNECT_RESP:
		p = &CmdConnectRspPkt{}
	case CMD_TERMINATE:
		p = &CmdTerminateReqPkt{}
	case CMD_TERMINATE_RESP:
		p = &CmdTerminateRspPkt{}
	case CMD_OPEN_CHANNEL:
		p = &CmdOpenChannelReqPkt{}
	case CMD_OPEN_CHANNEL_RESP:
		p = &CmdOpenChannelRspPkt{}
	case CMD_UMBRELLA_IN:
		p = &CmdUmbrellaInReqPkt{}
	case CMD_UMBRELLA_IN_RESP:
		p = &CmdUmbrellaInRspPkt{}
	case CMD_UMBRELLA_OUT:
		p = &CmdUmbrellaOutReqPkt{}
	case CMD_UMBRELLA_OUT_RESP:
		p = &CmdUmbrellaOutRspPkt{}
	case CMD_ACTIVE_TEST:
		p = &CmdActiveTestReqPkt{}
	case CMD_ACTIVE_TEST_RESP:
		p = &CmdActiveTestRspPkt{}

	default:
		p = nil
		return nil, ErrCommandIdNotSupported
	}

	err = p.Unpack(leftData)
	if err != nil {
		return nil, err
	}
	log.Printf("received packer: %v", p)
	return p, nil
}
