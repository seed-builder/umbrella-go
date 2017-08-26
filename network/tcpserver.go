package network

import (
	"time"
	"log"
	"net"
	"errors"
	"umbrella/utilities"
)

// errors for cmpp server
var (
	ErrEmptyServerAddr = errors.New("equipment server listen: empty server addr")
	ErrNoHandlers      = errors.New("equipment server: no connection handler")
	ErrUnsupportedPkt  = errors.New("equipment server read packet: receive a unsupported pkt")
)


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
				utilities.SysLog.Errorf("客户端接入错误：%v; retrying in %v", e, tempDelay)
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

		utilities.SysLog.Infof("收到客户端接入： %v", c.rw.RemoteAddr())
		go c.Serve()
	}
}

func (srv *TcpServer) newConn(rwc net.Conn) (c *Conn, err error) {
	c = NewConn(srv, rwc, V10)
	c.server = srv
	c.rw = rwc
	c.SetState(CONN_CONNECTED)
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

