package umbrella

import (
	"time"
	"io"
	"log"
	"os"
	"umbrella/network"
	"errors"
	"umbrella/models"
	"fmt"
	"umbrella/utilities"
	//"strconv"
	"strconv"
)

//EquipmentService is 单台设备管理服务，
type EquipmentService struct {
	EquipmentConns map[string]*network.Conn
	WaitingHire map[string]uint
	Requests map[uint8]chan int32
	Redoes map[uint8]*network.Packer
}

// ListenAndServe listens on the TCP network address addr
// and then calls Serve with handler to handle requests.
func (es *EquipmentService) ListenAndServe(addr string, ver  network.Type, t time.Duration, n int32, logWriter io.Writer) error {
	if addr == "" {
		return network.ErrEmptyServerAddr
	}

	//if handlers == nil {
	//	return network.ErrNoHandlers
	//}
	handlers := []network.Handler{
		network.HandlerFunc(es.HandleConnect),
		network.HandlerFunc(es.HandleUmbrellaIn),
		network.HandlerFunc(es.HandleUmbrellaOutRsp),
		network.HandlerFunc(es.HandleCmdIllegalRsp),
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
		ErrorLog: log.New(logWriter, "equipment server: ", log.LstdFlags)}
	return server.ListenAndServe()
}

func (es *EquipmentService) RegisterConn(equipmentSn string, conn *network.Conn)  {
	es.EquipmentConns[equipmentSn] = conn
}

func (es *EquipmentService) OpenChannel(equipmentSn string) (channelNum uint8, seqId uint8, err error) {
	conn, ok := es.EquipmentConns[equipmentSn]
	if ok {
		channelNum := conn.Equipment.ChooseChannel()
		req := &network.CmdUmbrellaOutReqPkt{}
		req.ChannelNum = channelNum
		seqId := <- conn.SeqId
		err := conn.SendPkt(req, seqId)

		if err != nil {
			return 0, 0, err
		} else {
			//重发
			go func() {
				time.Sleep(15 * time.Second)
				_, ok := es.Requests[seqId]
				if ok {
					log.Println("resend  OpenChannel request pkt !")
					conn.SendPkt(req, seqId)
				}
			}()

			_, ok := es.Requests[seqId]
			if !ok {
				es.Requests[seqId] = make(chan int32)
			}
			return channelNum , seqId, nil
		}
	} else {
		return 0, 0, errors.New("equipment is offline")
	}
}

func (es *EquipmentService) getKey(sn string, channelNum uint8) string {
	var k = fmt.Sprintf("%s%d", sn, channelNum)
	return k
}

func (es *EquipmentService) Close(){
	for sn, conn := range es.EquipmentConns {
		conn.Close()
		log.Println("close conn sn: ", sn)
	}
}

//HandleConnect
func (es *EquipmentService) HandleConnect(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	req, ok := p.Packer.(*network.CmdConnectReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}
	resp := r.Packer.(*network.CmdConnectRspPkt)
	resp.Status = utilities.RspStatusFail
	if req.EquipmentSn != "" {
		eq := models.Equipment{}
		eq.Query().First(&eq, "sn = ?", req.EquipmentSn)
		if eq.ID > 0 {
			eq.InitChannel()
			eq.Online()
			r.Packet.Conn.SetState( network.CONN_AUTHOK )
			r.Packet.Conn.SetEquipment(&eq)
			EquipmentSrv.RegisterConn(req.EquipmentSn, r.Packet.Conn)
			resp.Status = utilities.RspStatusSuccess
			l.Printf("connect success, sn: %s", req.EquipmentSn)
		}else{
			resp.Status = utilities.RspStatusEquipmentSnIllegal
			l.Printf("connect fail, sn: %s", req.EquipmentSn)
		}
	}
	return true, nil
}

//handleUmbrellaIn: umbrella in channel request
func (es *EquipmentService) HandleUmbrellaIn(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	req, ok := p.Packer.(*network.CmdUmbrellaInReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}
	l.Printf("handle the umbrella in request , %v", req)
	resp := r.Packer.(*network.CmdUmbrellaInRspPkt)
	if r.Packet.Conn.State != network.CONN_AUTHOK {
		resp.Status = utilities.RspStatusNeedAuth
		return false, nil
	}else{
		umbrella := models.Umbrella{}
		sn := strconv.Itoa(int(req.UmbrellaSn))
		resp.Status = umbrella.InEquipment(r.Packet.Conn.Equipment, sn, req.ChannelNum)
		return true, nil
	}
}

//HandleOpenChannelRsp
func (es *EquipmentService) HandleUmbrellaOutRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	rsp, ok := p.Packer.(*network.CmdUmbrellaOutRspPkt)
	if ok {
		c, o := es.Requests[rsp.SeqId]
		if o {
			log.Println("HandleUmbrellaOutRsp UmbrellaSn = ", rsp.UmbrellaSn)
			if rsp.Status == 1 {
				c <- rsp.UmbrellaSn
			}else{
				close(c)
			}
			delete( es.Requests, rsp.SeqId )
		}
	}
	return true, nil
}

//HandleCmdIllegalRsp
func (es *EquipmentService) HandleCmdIllegalRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	resp, ok := r.Packer.(*network.CmdIllegalRspPkt)
	if ok {
		resp.Status = utilities.RspStatusCmdIllegal
	}
	return false, nil
}


var EquipmentSrv *EquipmentService

func init()  {
	EquipmentSrv = &EquipmentService{
		EquipmentConns: make(map[string]*network.Conn),
		Requests: make(map[uint8]chan int32),
	}
}



