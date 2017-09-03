package main

import (
	"time"
	"sync"
	"umbrella/network"
	"umbrella/utilities"
	//"math/rand"
	//"encoding/hex"
)

const (
	connectTimeout time.Duration = time.Second * 2
)

var umbrellaSns = []string{
	"17617E62",
	"B7A87962",
	"87299762",
	"67C37462",
	"F7979162",
	"A7FB8262",
	"E7C27562",
	"67B45112",
}

func startAClient(idx int, sn string) {
	c := network.NewClient(0x10)
	defer wg.Done()
	defer c.Disconnect()
	//119.23.214.176, 39.108.180.41
	err := c.Connect(":7777", sn, connectTimeout)
	if err != nil {
		utilities.SysLog.Errorf("client %d: connect error: %s.", idx, err)
		return
	}
	utilities.SysLog.Infof("client %d: connect and auth ok", idx)

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			//p :=  &network.CmdUmbrellaInReqPkt{}
			//p.ChannelNum = uint8(rand.Intn(5))
			//var i = rand.Intn(5)
			//sn := umbrellaSns[i]
			//utilities.SysLog.Infof("client %d: prepare to send a umbrella in request  SN: %s,  ChannelNum: %d.", idx, sn, p.ChannelNum)
			//p.UmbrellaSn, _ = hex.DecodeString(sn)
			//err = c.SendReqPkt(p)
			//utilities.SysLog.Infof("client %d: send a umbrella in request : %v.", idx, p)
			//if err != nil {
			//	utilities.SysLog.Infof("client %d: send a umbrella in request error: %s.", idx, err)
			//} else {
			//	utilities.SysLog.Infof("client %d: send a umbrella in request ok", idx)
			//}
			break
		default:
		}

		// recv packets
		ps, err := c.RecvAndUnpackPkt(time.Second * 10)
		if err != nil {
			//utilities.SysLog.Infof("client %d: client read and unpack pkt error: %s.", idx, err)
			//break
			continue
		}
		i := ps[0]
		switch p := i.(type) {

		case *network.CmdActiveTestReqPkt:
			utilities.SysLog.Infof("client %d: receive a network active request: %v.", idx, p)
			rsp := &network.CmdActiveTestRspPkt{}
			err := c.SendRspPkt(rsp, p.SeqId)
			if err != nil {
				utilities.SysLog.Infof("client %d: send network active response error: %s.", idx, err)
				break
			}else{
				utilities.SysLog.Infof("client %d: send network active response success.", idx)
			}

		case *network.CmdUmbrellaInRspPkt:
			utilities.SysLog.Infof("client %d: receive a network umbrella in response: %v.", idx, p)

		case *network.CmdUmbrellaOutReqPkt:
			utilities.SysLog.Infof("client %d: receive a network open channel request: %v.", idx, p)
			rsp := &network.CmdUmbrellaOutRspPkt{
				Status: 1,
				UmbrellaSn: []byte{0x88, 0x04, 0xe3, 0x84},
			}
			time.Sleep(6*time.Second)
			err := c.SendRspPkt(rsp, p.SeqId)
			if err != nil {
				utilities.SysLog.Infof("client %d: send network open channel  response error: %s.", idx, err)
				break
			}else{
				utilities.SysLog.Infof("client %d: send network open channel  response success.", idx)
			}

		}
	}
}

var wg sync.WaitGroup

func main() {
	utilities.SysLog.Info("Client example start!")
	sn := []string{
		"E198402mqvw",
		"M201307olje",
		"E197005nlji",
		"M198505jboc",
		"M200808bryy",
		"M201506yqbu",
		"E200109tfjm",
		"E200702tlaj",
		"M200705kxra",
		"E198408egtb",
	}
	//for i := 0; i < 1; i++ {
		wg.Add(1)
		go startAClient(1, sn[0])
	//}
	wg.Wait()
	utilities.SysLog.Info("Client example ends!")
}
