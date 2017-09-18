package network

import (
	"log"
	"time"
)

type Packet struct {
	Packer
	*Conn
}

type Response struct {
	*Packet
	Packer
	SeqId uint8
	Status uint8
}

type Handler interface {
	ServeHandle(*Response, *Packet, *log.Logger) (bool, error)
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as Cmpp handlers.  If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler object that calls f.
//
// The first return value indicates whether to invoke next handler in
// the chain of handlers.
//
// The second return value shows the error returned from the handler. And
// if it is non-nil, server will close the client connection
// after sending back the corresponding response.
type HandlerFunc func(*Response, *Packet, *log.Logger) (bool, error)

// ServeHTTP calls f(r, p).
func (f HandlerFunc) ServeHandle(r *Response, p *Packet, l *log.Logger) (bool, error) {
	return f(r, p, l)
}

type Cmd struct {
	packer Packer
	sended time.Time
}

