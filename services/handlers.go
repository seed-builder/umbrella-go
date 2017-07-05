package services

import (
	"umbrella/network"
	"log"
	"umbrella/models"
)

func handleConnect(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
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
		eq.Query().First(&eq, "sn", req.EquipmentSn)
		if eq.ID > 0 {
			resp.Status = network.ConnectSuccess
			EquipmentSrv.RegisterConn(req.EquipmentSn, r.Packet.Conn)
		}else{
			resp.Status = network.ConnectWrongSn
		}
	}
	return true, nil
}

//handleUmbrellaIn: umbrella in channel request
func handleUmbrellaIn(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	req, ok := p.Packer.(*network.CmdUmbrellaOutReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}
	resp := r.Packer.(*network.CmdUmbrellaOutRspPkt)
	resp.Reserved = 1
	go EquipmentSrv.UmbrellaIn(req.EquipmentSn, req.UmbrellaSn, req.ChannelNum)
	return true, nil
}

//handleUmbrellaOut: umbrella out channel request
func handleUmbrellaOut(r *network.Response, p *network.Packet, l *log.Logger) (bool, error){
	req, ok := p.Packer.(*network.CmdUmbrellaInReqPkt)
	if !ok {
		// not a connect request, ignore it,
		// go on to next handler
		return true, nil
	}
	resp := r.Packer.(*network.CmdUmbrellaInRspPkt)
	resp.Reserved = 1
	go EquipmentSrv.UmbrellaOut(req.EquipmentSn, req.UmbrellaSn, req.ChannelNum)
	return true, nil
}