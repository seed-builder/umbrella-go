package network

import (
	"net"
	"errors"
	"time"
	"umbrella/models"
	"log"
)

type State uint8
type Type int8


const (
	V10 Type = 0x10
	CMDHEAD byte = 0xAA
	CMDFOOT byte = 0x55
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
	net.Conn
	State State
	Typ   Type

	// for SeqId generator goroutine
	SeqId <-chan uint8
	done  chan<- struct{}

	Equipment *models.Equipment
}

func newSeqIdGenerator() (<-chan uint8, chan<- struct{}) {
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
		if c.Equipment != nil {
			c.Equipment.Offline()
		}
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
func (c *Conn) SendPkt(packet Packer, seqId uint8) error {
	if c.State == CONN_CLOSED {
		return ErrConnIsClosed
	}
	buf := make([]byte, 0)
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

	_, err = c.Conn.Write(buf) //block write

	if err != nil {
		return err
	}

	return nil
}

const (
	defaultReadBufferSize = 64
)

// RecvAndUnpackPkt receives CMD byte stream, and unpack it to some CMD packet structure.
func (c *Conn) RecvAndUnpackPkt(timeout time.Duration) ([]Packer, error) {
	if c.State == CONN_CLOSED {
		return nil, ErrConnIsClosed
	}

	if timeout != 0 {
		t := time.Now().Add(timeout)
		c.SetReadDeadline(t)
		defer c.SetReadDeadline(noDeadline)
	}

	leftData := make([]byte, defaultReadBufferSize)
	length, err := c.Conn.Read(leftData) //io.ReadFull(c.Conn, leftData)
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
	var head ,foot , complete int
	for i := 0 ; i < len; i ++ {
		d := data[i]
		if d == CMDHEAD {
			head = i
			complete ++
		}
		if d == CMDFOOT {
			foot = i
			complete ++
		}
		if complete == 2 && foot > head {
			result = append(result, data[head+1:foot])
			head=0
			foot=0
			complete = 0
		}
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