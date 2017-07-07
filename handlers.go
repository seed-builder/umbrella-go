package umbrella

import (
	"umbrella/network"
	"log"
	"umbrella/models"
)

func HandleConnect(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
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
func HandleUmbrellaIn(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
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

//handleUmbrellaOut: umbrella out channel request
func HandleUmbrellaOut(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
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
	resp.Status = umbrella.OutEquipment(r.Packet.Conn.Equipment, req.UmbrellaSn, req.ChannelNum)
	return true, nil
}