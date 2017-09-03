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
)

//EquipmentService is 单台设备管理服务，
type EquipmentService struct {
	EquipmentConns map[string]*network.Conn
	Requests map[uint8]chan string
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
	server := &network.TcpServer{
		Addr: addr,
		Handler: handler,
		Typ: ver,
		T: t,
		N: n,
	}
	return server.ListenAndServe()
}

func (es *EquipmentService) RegisterConn(equipmentSn string, conn *network.Conn)  {
	es.EquipmentConns[equipmentSn] = conn
}

func (es *EquipmentService) OpenChannel(equipmentSn string, channelNum uint8) (uint8, uint8, error) {
	conn, ok := es.EquipmentConns[equipmentSn]
	if ok && conn.State > 0  {
		if channelNum == 0 {
			channelNum = conn.Equipment.ChooseChannel()
		}
		var seqId uint8
		req := &network.CmdUmbrellaOutReqPkt{}
		req.ChannelNum = channelNum
		seqId = <- conn.SeqId
		utilities.SysLog.Infof("发送设备【%s】开启通道【%d】 命令 序列号【%d】!", equipmentSn, channelNum, seqId)
		err := conn.SendPkt(req, seqId)

		if err != nil {
			return 0, 0, err
		} else {
			//重发
			go func() {
				time.Sleep(time.Duration(utilities.SysConfig.TcpResendInterval) * time.Second)
				_, ok := es.Requests[seqId]
				if ok {
					utilities.SysLog.Infof("重发设备【%s】开启通道【%d】 命令 序列号【%d】!", equipmentSn, channelNum, seqId)
					conn.SendPkt(req, seqId)
				}
			}()

			_, ok := es.Requests[seqId]
			if !ok {
				es.Requests[seqId] = make(chan string)
			}
			return channelNum , seqId, nil
		}
	} else {
		utilities.SysLog.Infof("设备【%s】离线无法发送命令!", equipmentSn)
		return 0, 0, errors.New("设备离线，无法发送命令")
	}
}

func (es *EquipmentService) getKey(sn string, channelNum uint8) string {
	var k = fmt.Sprintf("%s%d", sn, channelNum)
	return k
}

func (es *EquipmentService) Close(){
	utilities.SysLog.Notice("正常关闭服务")
	for sn, conn := range es.EquipmentConns {
		conn.Close()
		utilities.SysLog.Noticef("正常关闭设备【%s】连接", sn)
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
	utilities.SysLog.Infof("收到设备登陆命令, 设备【%s】 序列号【%d】", req.EquipmentSn, req.SeqId)
	resp := r.Packer.(*network.CmdConnectRspPkt)
	resp.Status = utilities.RspStatusFail
	if req.EquipmentSn != "" {
		eq := models.Equipment{}
		eq.Query().First(&eq, "sn = ?", req.EquipmentSn)
		if eq.ID > 0 {
			eq.InitChannel()
			eq.Online(r.Conn.Ip)
			r.Packet.Conn.SetState( network.CONN_AUTHOK )
			r.Packet.Conn.SetEquipment(&eq)
			es.RegisterConn(req.EquipmentSn, r.Packet.Conn)
			resp.Status = utilities.RspStatusSuccess
			utilities.SysLog.Infof("设备登陆成功, 设备【%s】 序列号【%d】", req.EquipmentSn, req.SeqId)
		}else{
			resp.Status = utilities.RspStatusEquipmentSnIllegal
			utilities.SysLog.Warningf("设备登陆失败, 设备【%s】 序列号【%d】", req.EquipmentSn, req.SeqId)
		}
	}
	return true, nil
}

//handleUmbrellaIn: umbrella in channel request
func (es *EquipmentService) HandleUmbrellaIn(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*network.CmdUmbrellaInReqPkt)
	if !ok {
		return true, nil
	}
	utilities.SysLog.Infof("收到还伞命令, 设备【%s】,通道【%d】,伞编号【%X】, 序列号【%d】", r.Equipment.Sn,
		req.ChannelNum, req.UmbrellaSn, req.SeqId)
	resp := r.Packer.(*network.CmdUmbrellaInRspPkt)

	umbrella := models.Umbrella{}
	sn := fmt.Sprintf("%X", req.UmbrellaSn)
	resp.Status = umbrella.InEquipment(r.Equipment, sn, req.ChannelNum)
	return true, nil

}

//HandleOpenChannelRsp
func (es *EquipmentService) HandleUmbrellaOutRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	rsp, ok := p.Packer.(*network.CmdUmbrellaOutRspPkt)
	if ok {
		c, o := es.Requests[rsp.SeqId]
		if o {
			utilities.SysLog.Infof("收到设备出伞反馈，设备【%s】,伞编号【%X】,状态【%s】,序列号【%d】",r.Equipment.Sn, rsp.UmbrellaSn, utilities.RspStatusDesc(rsp.Status), rsp.SeqId)
			if rsp.Status == 1 {
				c <- fmt.Sprintf("%X", rsp.UmbrellaSn)
			}else{
				close(c)
			}
			delete( es.Requests, rsp.SeqId )
			utilities.SysLog.Infof("删除等待序列，伞编号【%X】， 序列号【%d】", rsp.UmbrellaSn, rsp.SeqId)
		}
	}
	return true, nil
}

//HandleCmdIllegalRsp
func (es *EquipmentService) HandleCmdIllegalRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	resp, ok := r.Packer.(*network.CmdIllegalRspPkt)
	if ok {
		utilities.SysLog.Warning("收到非法命令")
		resp.Status = utilities.RspStatusCmdIllegal
	}
	return false, nil
}

var EquipmentSrv *EquipmentService

func init()  {
	EquipmentSrv = &EquipmentService{
		EquipmentConns: make(map[string]*network.Conn),
		Requests: make(map[uint8]chan string),
	}
}



