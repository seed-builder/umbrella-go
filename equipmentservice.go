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
	"sync/atomic"
)

//EquipmentService is 单台设备管理服务，
type EquipmentService struct {
	EquipmentConns map[string]*network.Conn
	//Requests map[uint8]chan string
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
		network.HandlerFunc(es.HandleDataErr),
		network.HandlerFunc(es.HandleActiveTestRsp),
		network.HandlerFunc(es.HandleConnect),
		network.HandlerFunc(es.HandleUmbrellaIn),
		network.HandlerFunc(es.HandleUmbrellaOutRsp),
		network.HandlerFunc(es.HandleChannelInspectRsp),
		network.HandlerFunc(es.HandleChannelRescueRsp),
		network.HandlerFunc(es.HandleTakeUmbrellaRsp),
		network.HandlerFunc(es.HandleUmbrellaInspect),
		network.HandlerFunc(es.ClearCmdCache),
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
		if channelNum == 0 {
			return 0 , 0, nil
		}
		var seqId uint8
		req := &network.CmdTakeUmbrellaReqPkt{}
		req.Channel = channelNum
		seqId = <- conn.SeqId
		utilities.SysLog.Infof("发送设备【%s】开启通道【%d】取伞命令 序列号【%d】!", equipmentSn, channelNum, seqId)
		err := conn.SendPkt(req, seqId)

		if err != nil {
			return 0, 0, err
		} else {
			//重发
			go func() {
				time.Sleep(time.Duration(utilities.SysConfig.TcpResendInterval) * time.Second)
				_, ok := conn.UmbrellaRequests[seqId]
				if ok {
					utilities.SysLog.Infof("重发设备【%s】开启通道【%d】取伞命令 序列号【%d】!", equipmentSn, channelNum, seqId)
					conn.SendPkt(req, seqId)
				}
			}()
			//超时
			go func() {
				time.Sleep(time.Duration(utilities.SysConfig.TcpResendInterval * 2) * time.Second)
				c, ok :=  conn.UmbrellaRequests[seqId]
				if ok {
					utilities.SysLog.Warningf("设备【%s】开启通道【%d】取伞命令超时 序列号【%d】!", equipmentSn, channelNum, seqId)
					c <- network.UmbrellaRequest{Success:false, Err: "借伞超时, 请重新扫码"}
					delete( conn.UmbrellaRequests, seqId )
				}
			}()

			_, ok := conn.UmbrellaRequests[seqId]
			if !ok {
				conn.UmbrellaRequests[seqId] = make(chan network.UmbrellaRequest)
			}
			return channelNum , seqId, nil
		}
	} else {
		utilities.SysLog.Infof("设备【%s】离线无法发送命令!", equipmentSn)
		return 0, 0, errors.New("设备离线，无法发送命令")
	}
}

func (es *EquipmentService) SetChannel(equipmentSn string, channelNum uint8, valid bool) {
	conn, ok := es.EquipmentConns[equipmentSn]
	if ok && conn.State > 0 {
		conn.SetChannelValid(channelNum, valid)
		//if valid {
		//	conn.SetChannelStatus(channelNum, utilities.RspStatusChannelMiddle)
		//} else {
		//	conn.SetChannelStatus(channelNum, utilities.RspStatusChannelErrLock)
		//}

	}
}

//BorrowUmbrella 从设备借伞
func (es *EquipmentService) BorrowUmbrella(customerId uint, equipmentSn string, channelNum uint8) (uint8, uint8, error) {
	conn, ok := es.EquipmentConns[equipmentSn]
	if ok && conn.State > 0  {
		if conn.ChannelInspectStatus == 0 {
			return 0, 0, errors.New("设备正在检查,请稍后")
		}
		if conn.RunStatus > network.RUN_STATUS_WAITING {
			return 0, 0, errors.New("设备忙，请重新扫码")
		}
		if channelNum == 0 {
			channelNum = conn.Equipment.ChooseChannel()
		}
		if channelNum == 0 {
			return 0 , 0, nil
		}
		//设置设备当前状态为：借伞
		conn.RunStatus = network.RUN_STATUS_BORROWING

		var seqId uint8
		req := &network.CmdTakeUmbrellaReqPkt{}
		req.Channel = channelNum
		seqId = <- conn.SeqId
		utilities.SysLog.Infof("发送设备【%s】开启通道【%d】取伞命令 序列号【%d】!", equipmentSn, channelNum, seqId)
		err := conn.SendPkt(req, seqId)

		if err != nil {
			return 0, 0, err
		} else {
			//重发
			//go func() {
			//	time.Sleep(time.Duration(utilities.SysConfig.TcpResendInterval) * time.Second)
			//	_, ok := conn.UmbrellaRequests[seqId]
			//	if ok {
			//		utilities.SysLog.Infof("重发设备【%s】开启通道【%d】取伞命令 序列号【%d】!", equipmentSn, channelNum, seqId)
			//		conn.SendPkt(req, seqId)
			//	}
			//}()
			//超时
			go func() {
				time.Sleep(time.Duration(utilities.SysConfig.TcpResendInterval * 1) * time.Second)
				c, ok :=  conn.UmbrellaRequests[seqId]
				if ok {
					utilities.SysLog.Warningf("设备【%s】开启通道【%d】取伞命令超时 序列号【%d】!", equipmentSn, channelNum, seqId)
					c <- network.UmbrellaRequest{Success:false, Err: "借伞超时, 请重新扫码"}
					delete( conn.UmbrellaRequests, seqId )
				}
			}()
			//记录借伞人ID
			//conn.Borrowers[seqId] = customerId
			_, ok := conn.UmbrellaRequests[seqId]
			if !ok {
				conn.UmbrellaRequests[seqId] = make(chan network.UmbrellaRequest)
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

func (es *EquipmentService) HandleDataErr(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	if r.Packer != nil && r.Status == utilities.RspStatusDataErr {
		return false, nil
	}
	return true, nil
}

func (es *EquipmentService) HandleActiveTestRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	rsp, ok := r.Packet.Packer.(*network.CmdActiveTestRspPkt)
	if ok {
		atomic.AddInt32(&p.Conn.Counter, -1)
		if r.Conn.State == network.CONN_AUTHOK {
			r.Conn.Equipment.Online(r.Conn.Ip)
		}
		if rsp.Channel > 0 && rsp.Status == utilities.RspStatusChannelErrLock {
			es.RescueChannel(r, p, rsp.Channel)
		}
	}
	return  true, nil
}

//HandleConnect 设备登陆服务器
func (es *EquipmentService) HandleConnect(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	req, ok := p.Packer.(*network.CmdConnectReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}
	utilities.SysLog.Noticef("收到设备【%s】登陆命令, 序列号【%d】", req.EquipmentSn, req.SeqId)
	resp := r.Packer.(*network.CmdConnectRspPkt)
	resp.Status = utilities.RspStatusUserWrong
	if req.EquipmentSn != "" {
		eq := &models.Equipment{}
		eq.Query().First(&eq, "sn = ?", req.EquipmentSn)
		if eq.ID > 0 {
			eq.ServerHttpBase = utilities.SysConfig.HttpBaseUrl
			eq.InitChannel()
			eq.Online(r.Conn.Ip)
			r.Packet.Conn.SetState( network.CONN_AUTHOK )
			r.Packet.Conn.SetEquipment(eq)
			es.RegisterConn(req.EquipmentSn, r.Packet.Conn)
			resp.Status = utilities.RspStatusSuccess
			utilities.SysLog.Noticef("设备【%s】登陆成功, 序列号【%d】", req.EquipmentSn, req.SeqId)
		}else{
			resp.Status = utilities.RspStatusEquipmentSnIllegal
			utilities.SysLog.Warningf("设备【%s】登陆失败, 序列号【%d】", req.EquipmentSn, req.SeqId)
		}
	}
	return true, nil
}

//handleUmbrellaIn: 通道进伞
func (es *EquipmentService) HandleUmbrellaIn(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	req, ok := p.Packer.(*network.CmdUmbrellaInReqPkt)
	if !ok {
		return true, nil
	}
	utilities.SysLog.Noticef("收到设备【%s】还伞命令,通道【%d】,伞编号【%X】, 序列号【%d】", r.Equipment.Sn,
		req.Channel, req.UmbrellaSn, req.SeqId)
	resp := r.Packer.(*network.CmdUmbrellaInRspPkt)
	// 检查通道是否有效
	if r.Equipment.CheckValid(req.Channel) {
		//设置设备当前状态为：还伞
		p.Conn.RunStatus = network.RUN_STATUS_RESTORE

		umbrella := &models.Umbrella{}
		sn := fmt.Sprintf("%X", req.UmbrellaSn)
		resp.Status = umbrella.InEquipment(r.Equipment, sn, req.Channel)
		p.Equipment.SetChannelStatus(req.Channel, utilities.RspStatusChannelReturn)
		//设置设备当前状态为：等待
		p.Conn.RunStatus = network.RUN_STATUS_WAITING
	}else{
		utilities.SysLog.Noticef("检测设备【%s】通道【%d】【无效】无法还伞, 序列号【%d】", r.Equipment.Sn,
			req.Channel, req.SeqId)
		resp.Status = utilities.RspStatusChannelErrLock
	}
	return true, nil
}

//HandleTakeUmbrella: 通道取伞， 成功则发送出伞命令
func (es *EquipmentService) HandleTakeUmbrellaRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	rsp, ok := p.Packer.(*network.CmdTakeUmbrellaRspPkt)
	if ok {
		utilities.SysLog.Noticef("收到设备【%s】设备取伞反馈,伞编号【%X】,状态【%s】,序列号【%d】",r.Equipment.Sn, rsp.UmbrellaSn, utilities.RspStatusDesc(rsp.Status), rsp.SeqId)
		c, ok := r.Conn.UmbrellaRequests[rsp.SeqId] //es.Requests[rsp.SeqId]
		if rsp.Status == utilities.RspStatusSuccess {
			umbrella := &models.Umbrella{}
			sn := fmt.Sprintf("%X", rsp.UmbrellaSn)
			status := umbrella.Check(sn)
			utilities.SysLog.Noticef("设备【%s】取伞反馈,伞编号【%s】,查询出的该伞状态是【%s】,序列号【%d】",r.Equipment.Sn, sn, utilities.RspStatusDesc(status), rsp.SeqId)

			ur := network.UmbrellaRequest{}
			if status == utilities.RspStatusSuccess {
				ur.Success = true
				ur.Sn = sn
			}else{
				r.Packer = nil
				ur.Err = utilities.RspStatusDesc(status)
			}
			if ok {
				c <- ur
				delete(r.Conn.UmbrellaRequests, rsp.SeqId)
				utilities.SysLog.Infof("删除等待序列，设备【%s】, 伞编号【%X】， 序列号【%d】", r.Equipment.Sn, rsp.UmbrellaSn, rsp.SeqId)
			}
		}else{
			if ok {
				c <- network.UmbrellaRequest{Success:false, Err: utilities.RspStatusDesc(rsp.Status) }
				delete(r.Conn.UmbrellaRequests, rsp.SeqId)
			}
			r.Packer = nil
			return es.HandleException(r, p, r.SeqId, rsp.Status, rsp.Channel)
		}
	}
	return true, nil
}

//HandleUmbrellaOutRsp 通道出伞
func (es *EquipmentService) HandleUmbrellaOutRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error) {
	rsp, ok := p.Packer.(*network.CmdUmbrellaOutRspPkt)
	if ok {
		//c, o := es.Requests[rsp.SeqId]
		//if o {
		utilities.SysLog.Noticef("收到设备【%s】出伞反馈,伞编号【%X】,状态【%s】,序列号【%d】", r.Equipment.Sn, rsp.UmbrellaSn, utilities.RspStatusDesc(rsp.Status), rsp.SeqId)
		//if rsp.Status == utilities.RspStatusSuccess {
		//	c <- fmt.Sprintf("%X", rsp.UmbrellaSn)
		//}else{
		//	close(c)
		//	es.HandleException(r, p, r.SeqId, rsp.Status, rsp.Channel)
		//}
		//delete( es.Requests, rsp.SeqId )
		//utilities.SysLog.Infof("删除等待序列，设备【%s】, 伞编号【%X】， 序列号【%d】", r.Equipment.Sn, rsp.UmbrellaSn, rsp.SeqId)
		//}
	}
	return true, nil
}

//HandleChannelInspectRsp: 设备通道检查反馈（服务端发起）
func (es *EquipmentService) HandleChannelInspectRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	rsp, ok := p.Packer.(*network.CmdChannelInspectRspPkt)
	if ok {
		utilities.SysLog.Noticef("收到设备【%s】通道检查反馈, 通道【%d】,状态【%s】,序列号【%d】",r.Equipment.Sn, rsp.Channel, utilities.RspStatusDesc(rsp.Status), rsp.SeqId)
		rescue :=  p.Conn.Equipment.SetChannelStatus(rsp.Channel, rsp.Status)
		if rescue {
			es.RescueChannel(r, p, rsp.Channel)
		}
		if rsp.Channel < p.Conn.Equipment.Channels {
			nextId := rsp.Channel + 1
			p.Conn.ChannelInspect(nextId)
		}else{
			utilities.SysLog.Noticef("设备【%s】通道检查完毕",r.Equipment.Sn)
			p.Conn.ChannelInspectStatus = 1
			p.Conn.RunStatus = network.RUN_STATUS_WAITING
		}
	}
	return true, nil
}

//HandleChannelInspectRsp: 设备通道救援反馈（服务端发起）
func (es *EquipmentService) HandleChannelRescueRsp(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	rsp, ok := p.Packer.(*network.CmdChannelRescueRspPkt)
	if ok {
		utilities.SysLog.Noticef("收到设备【%s】通道救援反馈, 通道【%d】,状态【%s】,序列号【%d】",r.Equipment.Sn, rsp.Channel, utilities.RspStatusDesc(rsp.Status), rsp.SeqId)
		rescue := p.Conn.Equipment.SetChannelStatus(rsp.Channel, rsp.Status)
		if rescue {
			//开启通道救援（恢复）命令
			seqId := <-p.Conn.SeqId
			r.Packer = &network.CmdChannelRescueReqPkt{
				CmdData: network.CmdData{
					SeqId:   seqId,
					Channel: rsp.Channel,
				},
			}
		}
	}
	return true, nil
}

//HandleUmbrellaInspect 设备发起的伞SN 检查
func (es *EquipmentService) HandleUmbrellaInspect(r *network.Response, p *network.Packet, l *log.Logger) (bool, error)  {
	req, ok := p.Packer.(*network.CmdUmbrellaInspectReqPkt)
	if ok {
		resp := r.Packer.(*network.CmdUmbrellaInspectRspPkt)
		utilities.SysLog.Noticef("收到设备【%s】伞SN检查命令, 通道【%d】,伞编号【%X】, 序列号【%d】", r.Equipment.Sn,
			req.Channel, req.UmbrellaSn, req.SeqId)
		// 检查通道是否有效
		if r.Equipment.CheckValid(req.Channel) {
			//设置设备当前状态为：还伞
			p.Conn.RunStatus = network.RUN_STATUS_RESTORE

			umbrella := &models.Umbrella{}
			sn := fmt.Sprintf("%X", req.UmbrellaSn)
			status := umbrella.InEquipment(r.Equipment, sn, req.Channel)

			p.Equipment.SetChannelStatus(req.Channel, utilities.RspStatusChannelBorrow)
			//设置设备当前状态为：等待
			p.Conn.RunStatus = network.RUN_STATUS_WAITING

			if status == utilities.RspStatusSuccess {
				resp.Status = utilities.RspStatusUmbrellaReturned
			} else {
				resp.Status = status
			}
		}else{
			utilities.SysLog.Noticef("检测设备【%s】通道【%d】【无效】无法还伞, 序列号【%d】", r.Equipment.Sn,
				req.Channel, req.SeqId)
			resp.Status = utilities.RspStatusChannelErrLock
		}
	}
	return true, nil
}

//ClearCmdCache 最后清理发送过的命令缓存
func (es *EquipmentService) ClearCmdCache(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	_, ok := p.SendCmdCache[r.SeqId]
	if ok {
		delete(p.SendCmdCache, r.SeqId)
	}
	return true, nil
}

//CatchException 异常状态处理
func (es *EquipmentService) HandleException(r *network.Response, p *network.Packet, seqId uint8, status uint8, channel uint8)  (bool, error) {
	desc := ""
	msg := &models.Message{}
	switch status {
	case utilities.RspStatusChannelBusy:
		desc = "通道忙（有指令未完成）"
	case utilities.RspStatusTimeout:
		fallthrough
	case utilities.RspStatusChannelErrLock:
		desc = "通道锁异常/超时"
		rescue := p.Equipment.SetChannelStatus(channel, status)
		if rescue {
			//开启通道救援（恢复）命令
			es.RescueChannel(r, p, channel)
		}
		desc = ""
	case utilities.RspStatusGprsErr:
		desc = "网络错误"
	case utilities.RspStatusUnknowError:
		desc = "未知错误"
	case utilities.RspStatusDataErr:
		desc = "数据错"
	case utilities.RspStatusChannelTimeout:
		desc = "通道超时"
	case utilities.RspStatusChannelMiddle:
		desc = "通道锁状态-中间"
	case utilities.RspStatusChannelBorrow:
		desc = "通道锁状态-借伞"
	case utilities.RspStatusChannelReturn:
		desc = "通道锁状态-还伞"
	case utilities.RspStatusChannelErr:
		desc = "通道命令不支持"
	case utilities.RspStatusNotMatch:
		desc = "通道和命令不匹配"
	case utilities.RspStatusChannelErrSN:
		desc = "通道伞SN不匹配"

	}
	if desc != ""{
		msg.AddEquipmentError(p.Equipment.Sn, p.Equipment.ID, p.Equipment.SiteId, desc)
	}
	return true, nil
}

func (es *EquipmentService) RescueChannel(r *network.Response,  p *network.Packet, channel uint8){
	//seqId := <- p.Conn.SeqId
	//r.Packer = &network.CmdChannelRescueReqPkt{
	//	CmdData: network.CmdData{
	//		SeqId:   seqId,
	//		Channel: channel,
	//	},
	//}
}

var EquipmentSrv *EquipmentService

func init()  {
	EquipmentSrv = &EquipmentService{
		EquipmentConns: make(map[string]*network.Conn),
		//Requests: make(map[uint8]chan string),
	}
}



