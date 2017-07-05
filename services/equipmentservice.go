package services

import (
	"umbrella/models"
	"time"
	"io"
	"log"
	"os"
	"umbrella/network"
	"github.com/kardianos/govendor/context"
)


//EquipmentService is 单台设备管理服务，
type EquipmentService struct {
	EquipmentConns map[string]*network.Conn
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests.
func (es *EquipmentService) ListenAndServe(addr string, ver  network.Type, t time.Duration, n int32, logWriter io.Writer, handlers ... network.Handler) error {
	if addr == "" {
		return network.ErrEmptyServerAddr
	}

	if handlers == nil {
		return network.ErrNoHandlers
	}

	var handler network.Handler
	handler = network.HandlerFunc(func(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
		for _, h := range handlers {
			next, err := h.ServeHandle(r, p, l)
			if err != nil || !next {
				return next, err
			}
		}
		return false, nil
	})

	if logWriter == nil {
		logWriter = os.Stderr
	}
	server := &network.TcpServer{Addr: addr, Handler: handler, Typ: ver,
		T: t, N: n,
		ErrorLog: log.New(logWriter, "cmd server: ", log.LstdFlags)}
	return server.ListenAndServe()
}

func (es *EquipmentService) UmbrellaIn(equipmentSn string, umbrellaSn string, channelNum uint8){

}

func (es *EquipmentService) UmbrellaOut(equipmentSn string, umbrellaSn string, channelNum uint8){

}

func (es *EquipmentService) RegisterConn(equipmentSn string, conn *network.Conn)  {
	es.EquipmentConns[equipmentSn] = conn
}

var EquipmentSrv *EquipmentService

func init()  {
	EquipmentSrv = &EquipmentService{
		EquipmentConns: make(map[string]*network.Conn),
	}
}



