package network

import (
	"time"
	"errors"
	"net"
	"strconv"
)

var ErrNotCompleted = errors.New("data not being handled completed")
var ErrRespNotMatch = errors.New("the response is not matched with the request")

// Client stands for one client-side instance, just like a session.
// It may connect to the server, send & recv cmpp packets and terminate the connection.
type Client struct {
	conn *Conn
	typ  Type
}

// New establishes a new cmpp client.
func NewClient(typ Type) *Client {
	return &Client{
		typ: typ,
	}
}

// Connect connect to the cmpp server in block mode.
// It sends login packet, receive and parse connect response packet.
func (cli *Client) Connect(servAddr, sn string, timeout time.Duration) error {
	var err error
	conn, err := net.DialTimeout("tcp", servAddr, timeout)
	if err != nil {
		return err
	}
	cli.conn = NewConn(conn, cli.typ)
	defer func() {
		if err != nil {
			cli.conn.Close()
		}
	}()
	cli.conn.SetState(CONN_CONNECTED)

	// Login to the server.
	req := &CmdConnectReqPkt{
		EquipmentSn: sn,
	}

	err = cli.SendReqPkt(req)
	if err != nil {
		return err
	}

	p, err := cli.conn.RecvAndUnpackPkt(0)
	if err != nil {
		return err
	}

	var ok bool
	var status uint8

	var rsp *CmdConnectRspPkt
	rsp, ok = p.(*CmdConnectRspPkt)
	status = rsp.Status


	if !ok {
		err = ErrRespNotMatch
		return err
	}

	if status != 1 {
		//err = ConnRspStatusErrMap[status]
		return errors.New("Conn Rsp StatusErr status: " + strconv.Itoa(int(status)))
	}

	cli.conn.SetState(CONN_AUTHOK)
	return nil
}

func (cli *Client) Disconnect() {
	cli.conn.Close()
}

// SendReqPkt pack the cmpp request packet structure and send it to the other peer.
func (cli *Client) SendReqPkt(packet Packer) error {
	return cli.conn.SendPkt(packet, <-cli.conn.SeqId)
}

// SendRspPkt pack the cmpp response packet structure and send it to the other peer.
func (cli *Client) SendRspPkt(packet Packer, seqId uint32) error {
	return cli.conn.SendPkt(packet, seqId)
}

// RecvAndUnpackPkt receives cmpp byte stream, and unpack it to some cmpp packet structure.
func (cli *Client) RecvAndUnpackPkt(timeout time.Duration) (interface{}, error) {
	return cli.conn.RecvAndUnpackPkt(timeout)
}
