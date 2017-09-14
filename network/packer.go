package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"errors"
)

// Common errors.
var ErrMethodParamsInvalid = errors.New("params passed to method is invalid")


type Packer interface {
	Pack(seqId uint8) ([]byte, error)
	Unpack(data []byte) error
}

// OpError is the error type usually returned by functions in the cmpppacket
// package. It describes the operation and the error which the operation caused.
type OpError struct {
	// err is the error that occurred during the operation.
	// it is the origin error.
	err error

	// op is the operation which caused the error, such as
	// some "read" or "write" in packetWriter or packetReader.
	op string
}

func NewOpError(e error, op string) *OpError {
	return &OpError{
		err: e,
		op:  op,
	}
}

func (e *OpError) Error() string {
	if e.err == nil {
		return "<nil>"
	}
	return e.op + " error: " + e.err.Error()
}

func (e *OpError) Cause() error {
	return e.err
}

func (e *OpError) Op() string {
	return e.op
}

type packetWriter struct {
	wb  *bytes.Buffer
	err *OpError
}

func newPacketWriter(initSize uint32) *packetWriter {
	buf := make([]byte, 0, initSize)
	return &packetWriter{
		wb: bytes.NewBuffer(buf),
	}
}

// Bytes returns a slice of the contents of the inner buffer;
// If the caller changes the contents of the
// returned slice, the contents of the buffer will change provided there
// are no intervening method calls on the Buffer.
func (w *packetWriter) Bytes() ([]byte, error) {
	if w.err != nil {
		return nil, w.err
	}
	len := w.wb.Len()
	return (w.wb.Bytes())[:len], nil
}

// WriteInt appends the byte of b to the inner buffer, growing the buffer as
// needed.
func (w *packetWriter) WriteByte(b byte) {
	if w.err != nil {
		return
	}
	err := w.wb.WriteByte(b)
	if err != nil {
		w.err = NewOpError(err,
			fmt.Sprintf("packetWriter.WriteByte writes: %x", b))
		return
	}
}

// WriteInt appends the byte of b to the inner buffer, growing the buffer as
// needed.
func (w *packetWriter) WriteBytes(b []byte) {
	if w.err != nil {
		return
	}
	_, err := w.wb.Write(b)
	if err != nil {
		w.err = NewOpError(err,
			fmt.Sprintf("packetWriter.WriteBytes writes: %x", b))
		return
	}
}

// WriteFixedSizeString writes a string to buffer, if the length of s is less than size,
// Pad binary zero to the left bytes.
func (w *packetWriter) WriteFixedSizeString(s string, size int) {
	if w.err != nil {
		return
	}

	l1 := len(s)
	l2 := l1
	if l2 > 10 {
		l2 = 10
	}

	if l1 > size {
		w.err = NewOpError(ErrMethodParamsInvalid,
			fmt.Sprintf("packetWriter.WriteFixedSizeString writes: %s", s[0:l2]))
		return
	}

	w.WriteString(strings.Join([]string{s, string(make([]byte, size-l1))}, ""))
}

// WriteString appends the contents of s to the inner buffer, growing the buffer as
// needed.
func (w *packetWriter) WriteString(s string) {
	if w.err != nil {
		return
	}

	l1 := len(s)
	l2 := l1
	if l2 > 10 {
		l2 = 10
	}

	n, err := w.wb.WriteString(s)
	if err != nil {
		w.err = NewOpError(err,
			fmt.Sprintf("packetWriter.WriteString writes: %s", s[0:l2]))
		return
	}

	if n != l1 {
		w.err = NewOpError(fmt.Errorf("WriteString writes %d bytes, not equal to %d we expected", n, l1),
			fmt.Sprintf("packetWriter.WriteString writes: %s", s[0:l2]))
		return
	}
}

// WriteInt appends the content of data to the inner buffer in order, growing the buffer as
// needed.
func (w *packetWriter) WriteInt(order binary.ByteOrder, data interface{}) {
	if w.err != nil {
		return
	}

	err := binary.Write(w.wb, order, data)
	if err != nil {
		w.err = NewOpError(err,
			fmt.Sprintf("packetWriter.WriteInt writes: %#v", data))
		return
	}
}

const maxCStringSize = 160

type packetReader struct {
	rb   *bytes.Buffer
	err  *OpError
	cbuf [maxCStringSize]byte
}

func newPacketReader(data []byte) *packetReader {
	return &packetReader{
		rb: bytes.NewBuffer(data),
	}
}

// ReadByte reads and returns the next byte from the inner buffer.
// If no byte is available, it returns an OpError.
func (r *packetReader) ReadByte() byte {
	if r.err != nil {
		return 0
	}

	b, err := r.rb.ReadByte()
	if err != nil {
		r.err = NewOpError(err,
			"packetReader.ReadByte")
		return 0
	}
	return b
}

// Read reads structured binary data from r into data.
// Data must be a pointer to a fixed-size value or a slice
// of fixed-size values.
// Bytes read from r are decoded using the specified byte order
// and written to successive fields of the data.
func (r *packetReader) ReadInt(order binary.ByteOrder, data interface{}) {
	if r.err != nil {
		return
	}
	err := binary.Read(r.rb, order, data)
	if err != nil {
		r.err = NewOpError(err,
			"packetReader.ReadInt")
		return
	}
}

// ReadBytes reads the next len(s) bytes from the inner buffer to s.
// If the buffer has no data to return, an OpError would be stored in r.err.
func (r *packetReader) ReadBytes(s []byte) {
	if r.err != nil {
		return
	}

	n, err := r.rb.Read(s)
	if err != nil {
		r.err = NewOpError(err,
			"packetReader.ReadBytes")
		return
	}

	if n != len(s) {
		r.err = NewOpError(fmt.Errorf("ReadBytes reads %d bytes, not equal to %d we expected", n, len(s)),
			"packetWriter.ReadBytes")
		return
	}
}

// ReadCString read bytes from packerReader's inner buffer,
// it would trim the tail-zero byte and the bytes after that.
func (r *packetReader) ReadCString(length int) []byte {
	if r.err != nil {
		return nil
	}

	var tmp = r.cbuf[:length]
	n, err := r.rb.Read(tmp)
	if err != nil {
		r.err = NewOpError(err,
			"packetReader.ReadCString")
		return nil
	}

	if n != length {
		r.err = NewOpError(fmt.Errorf("ReadCString reads %d bytes, not equal to %d we expected", n, length),
			"packetWriter.ReadCString")
		return nil
	}

	i := bytes.IndexByte(tmp, 0)
	if i == -1 {
		return tmp
	} else {
		return tmp[:i]
	}
}

// Error return the inner err.
func (r *packetReader) Error() error {
	if r.err != nil {
		return r.err
	}
	return nil
}
