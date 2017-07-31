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
)

//EquipmentService is 单台设备管理服务，
type EquipmentService struct {
	EquipmentConns map[string]*network.Conn
	WaitingHire map[string]uint
	Requests map[uint8]chan struct{}
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
		network.HandlerFunc(es.HandleUmbrellaOut),
		network.HandlerFunc(es.HandleOpenChannelRsp),
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
		req := &network.CmdOpenChannelReqPkt{}
		req.ChannelNum = channelNum
		seqId := <- conn.SeqId
		err := conn.SendPkt(req, seqId)

		if err != nil {
			return 0, 0, err
		} else {
			//重发
			go func() {
				time.Sleep(5 * time.Second)
				_, ok := es.Requests[seqId]
				if ok {
					log.Println("resend  OpenChannel request pkt !")
					conn.SendPkt(req, seqId)
				}
			}()

			_, ok := es.Requests[seqId]
			if !ok {
				es.Requests[seqId] = make(chan struct{})
			}
			return channelNum , seqId, nil
		}
	} else {
		return 0, 0, errors.New("equipment is offline")
	}
}

func (es *EquipmentService) DoHire(hire_id uint) (success bool, err error) {
	hire := models.CustomerHire{}
	hire.Query().Preload("HireEquipment").First(&hire, hire_id)
	conn, ok := es.EquipmentConns[hire.HireEquipment.Sn]
	if ok {
		channelNum := conn.Equipment.ChooseChannel()
		req := &network.CmdOpenChannelReqPkt{}
		req.ChannelNum = channelNum
		err := conn.SendPkt(req, <- conn.SeqId)
		if err != nil {
			return false, err
		} else {
			k := es.getKey(hire.HireEquipment.Sn, channelNum)
			es.WaitingHire[k] = hire_id
			return true , nil
		}
	} else{
		return false, errors.New("equipment is offline")
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
	resp.Status = network.ConnectFail
	if req.EquipmentSn != "" {
		eq := models.Equipment{}
		eq.Query().First(&eq, "sn = ?", req.EquipmentSn)
		if eq.ID > 0 {
			eq.InitChannel()
			eq.Online()
			r.Packet.Conn.SetState( network.CONN_AUTHOK )
			r.Packet.Conn.SetEquipment(&eq)
			EquipmentSrv.RegisterConn(req.EquipmentSn, r.Packet.Conn)
			resp.Status = network.ConnectSuccess
			l.Printf("connect success, sn: %s", req.EquipmentSn)
		}else{
			resp.Status = network.ConnectWrongSn
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
	if r.Packet.Conn.State != network.CONN_AUTHOK {
		return false, network.ErrConnNeedAuth
	}
	l.Printf("handle the umbrella in request , %v", req)
	resp := r.Packer.(*network.CmdUmbrellaInRspPkt)
	umbrella := models.Umbrella{}
	resp.Status = umbrella.InEquipment(r.Packet.Conn.Equipment, req.UmbrellaSn, req.ChannelNum)
	return true, nil
}

//HandleOpenChannelRsp
func (es *EquipmentService) HandleOpenChannelRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	rsp, ok := p.Packer.(*network.CmdOpenChannelRspPkt)
	if ok {
		c, o := es.Requests[rsp.SeqId]
		if o {
			log.Println("close request seqid = ", rsp.SeqId)
			close(c)
			delete( es.Requests, rsp.SeqId )
		}
	}
	return true, nil
}

//handleUmbrellaOut: umbrella out channel request
func (es *EquipmentService) HandleUmbrellaOut(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	req, ok := p.Packer.(*network.CmdUmbrellaOutReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}
	if r.Packet.Conn.State != network.CONN_AUTHOK {
		return false, network.ErrConnNeedAuth
	}
	l.Printf("handle the umbrella out request , %v", req)
	resp := r.Packer.(*network.CmdUmbrellaOutRspPkt)
	umbrella := models.Umbrella{}
	k := es.getKey(r.Packet.Conn.Equipment.Sn, req.ChannelNum)
	hire_id, ok := es.WaitingHire[k]
	if ok {
		delete(es.WaitingHire, k)
	}
	resp.Status = umbrella.OutEquipment(r.Packet.Conn.Equipment, req.UmbrellaSn, req.ChannelNum, hire_id)
	return true, nil
}

var EquipmentSrv *EquipmentService

func init()  {
	EquipmentSrv = &EquipmentService{
		EquipmentConns: make(map[string]*network.Conn),
		Requests: make(map[uint8]chan struct{}),
	}
}



