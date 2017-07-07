package network

import (
	"time"
	"log"
	"net"
	"errors"
	"fmt"
	"sync/atomic"
)

// errors for cmpp server
var (
	ErrEmptyServerAddr = errors.New("equipment server listen: empty server addr")
	ErrNoHandlers      = errors.New("equipment server: no connection handler")
	ErrUnsupportedPkt  = errors.New("equipment server read packet: receive a unsupported pkt")
)

// A conn represents the server side of a Cmd connection.
type conn struct {
	*Conn
	server *TcpServer // the Server on which the connection arrived

	// for active test
	t       time.Duration // interval between two active tests
	n       int32         // continuous send times when no response back
	done    chan struct{}
	exceed  chan struct{}
	counter int32
}

// Serve a new connection.
func (c *conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			c.server.ErrorLog.Printf("panic serving %v: %v\n", c.Conn.RemoteAddr(), err)
		}
	}()

	defer c.close()

	// start a goroutine for sending active test.
	//startActiveTest(c)

	for {
		select {
		case <-c.exceed:
			return // close the connection.
		default:
		}

		r, err := c.readPacket()
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				continue
			}
			break
		}

		_, err = c.server.Handler.ServeHandle(r, r.Packet, c.server.ErrorLog)
		if err1 := c.finishPacket(r); err1 != nil {
			break
		}

		if err != nil {
			break
		}
	}
}

func (c *conn) readPacket() (*Response, error) {
	readTimeout := time.Second * 2
	i, err := c.Conn.RecvAndUnpackPkt(readTimeout)
	if err != nil {
		return nil, err
	}
	//typ := c.server.Typ

	var pkt *Packet
	var rsp *Response
	switch p := i.(type) {
	case *CmdConnectReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdConnectRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmd connect request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdOpenChannelReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &CmdOpenChannelRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmd open channel request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdOpenChannelRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}
		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmd open channel response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdUmbrellaInReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaInRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmd umbrella in request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdUmbrellaOutReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}
		rsp = &Response{
			Packet: pkt,
			Packer: &CmdUmbrellaOutRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmd umbrella out request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdActiveTestReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &CmdActiveTestRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmd active request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdActiveTestRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmd active response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdTerminateReqPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
			Packer: &CmdTerminateRspPkt{
				SeqId: p.SeqId,
			},
			SeqId: p.SeqId,
		}
		c.server.ErrorLog.Printf("receive a cmpp terminate request from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	case *CmdTerminateRspPkt:
		pkt = &Packet{
			Packer: p,
			Conn:   c.Conn,
		}

		rsp = &Response{
			Packet: pkt,
		}
		c.server.ErrorLog.Printf("receive a cmpp terminate response from %v[%d]\n",
			c.Conn.RemoteAddr(), p.SeqId)

	default:
		return nil, NewOpError(ErrUnsupportedPkt,
			fmt.Sprintf("readPacket: receive unsupported packet type: %#v", p))
	}
	return rsp, nil
}

// Close the connection.
func (c *conn) close() {
	p := &CmdTerminateReqPkt{}

	err := c.Conn.SendPkt(p, <-c.Conn.SeqId)
	if err != nil {
		c.server.ErrorLog.Printf("send cmd terminate request packet to %v error: %v\n", c.Conn.RemoteAddr(), err)
	}

	close(c.done)
	c.server.ErrorLog.Printf("close connection with %v!\n", c.Conn.RemoteAddr())
	c.Conn.Close()
}

func (c *conn) finishPacket(r *Response) error {
	if _, ok := r.Packet.Packer.(*CmdActiveTestRspPkt); ok {
		atomic.AddInt32(&c.counter, -1)
		return nil
	}

	if r.Packer == nil {
		// For response packet received, it need not
		// to send anything back.
		return nil
	}

	return c.Conn.SendPkt(r.Packer, r.SeqId)
}

func startActiveTest(c *conn) {
	exceed, done := make(chan struct{}), make(chan struct{})
	c.done = done
	c.exceed = exceed

	go func() {
		t := time.NewTicker(c.t)
		defer t.Stop()
		for {
			select {
			case <-done:
				// once conn close, the goroutine should exit
				return
			case <-t.C:
				// check whether c.counter exceeds
				if atomic.LoadInt32(&c.counter) >= c.n {
					c.server.ErrorLog.Printf("no client active test response returned from %v for %d times!",
						c.Conn.RemoteAddr(), c.n)
					exceed <- struct{}{}
					break
				}
				// send a active test packet to peer, increase the active test counter
				p := &CmdActiveTestReqPkt{}
				err := c.Conn.SendPkt(p, <-c.Conn.SeqId)
				if err != nil {
					c.server.ErrorLog.Printf("send cmd active test request to %v error: %v", c.Conn.RemoteAddr(), err)
				} else {
					atomic.AddInt32(&c.counter, 1)
				}
			}
		}
	}()
}

type TcpServer struct {
	Addr    string
	Handler Handler

	// protocol info
	Typ Type
	T   time.Duration // interval betwwen two active tests
	N   int32         // continuous send times when no response back

	// ErrorLog specifies an optional logger for errors accepting
	// connections and unexpected behavior from handlers.
	// If nil, logging goes to os.Stderr via the log package's
	// standard logger.
	ErrorLog *log.Logger
}

func (srv *TcpServer) ListenAndServe() error {
	if srv.Addr == "" {
		return ErrEmptyServerAddr
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	return srv.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

// Serve accepts incoming connections on the Listener l, creating a
// new service goroutine for each.  The service goroutines read requests and
// then call srv.Handler to reply to them.
func (srv *TcpServer) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.ErrorLog.Printf("accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}

		srv.ErrorLog.Printf("accept a connection from %v\n", c.Conn.RemoteAddr())
		go c.serve()
	}
}

func (srv *TcpServer) newConn(rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.server = srv
	c.Conn = NewConn(rwc, srv.Typ)
	c.Conn.SetState(CONN_CONNECTED)
	c.n = c.server.N
	c.t = c.server.T
	return c, nil
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away. the tcpKeepAliveListener's implementation is copied from
// http package.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(1 * time.Minute) // 1min
	return tc, nil
}

