package utilities

import (
	"bytes"
	"github.com/op/go-logging"
)

type IDBLog interface {
	NewLog(level int, content string) bool
}

type DBackend struct {
	Logger IDBLog
}


func (b *DBackend) Log(level logging.Level, calldepth int, rec *logging.Record) error {
	if level > logging.NOTICE{
		return nil
	}
	buf := &bytes.Buffer{}
	buf.Write([]byte(rec.Formatted(calldepth + 1)))
 	b.Logger.NewLog(int(level), buf.String())
	return nil
}