package main

import (
	"time"
	"log"
	"sync"
	"umbrella/network"
	"math/rand"
)

const (
	connectTimeout time.Duration = time.Second * 2
)

var umbrellaSns = []string{
	"S200402zhjd",
	"S197709awfg",
	"S199803lfta",
	"S200907qczp",
	"S200608qjri",
	"S200601lmkz",
	"S197206spup",
	"S199409grcr",
	"S201601puei",
	"S200007xrqi",
	"S200805xieh",
	"S197811pjxy",
	"S201608dtix",
	"S198607qpkz",
	"S198804gcdw",
	"S199601jvsg",
	"S199501rohv",
	"S197712rjlw",
	"S197608pgww",
	"S199811wpvw",
}

func startAClient(idx int, sn string) {
	c := network.NewClient(0x10)
	defer wg.Done()
	defer c.Disconnect()
	err := c.Connect(":7777", sn, connectTimeout)
	if err != nil {
		log.Printf("client %d: connect error: %s.", idx, err)
		return
	}
	log.Printf("client %d: connect and auth ok", idx)

	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-t.C:

			p :=  &network.CmdUmbrellaInReqPkt{}
			p.ChannelNum = 1
			var i = rand.Intn(19)
			p.UmbrellaSn = umbrellaSns[i]
			err = c.SendReqPkt(p)
			log.Printf("client %d: send a umbrella in request : %v.", idx, p)
			if err != nil {
				log.Printf("client %d: send a umbrella in request error: %s.", idx, err)
			} else {
				log.Printf("client %d: send a umbrella in request ok", idx)
			}
			break
		default:
		}

		// recv packets
		i, err := c.RecvAndUnpackPkt(time.Second * 10)
		if err != nil {
			//log.Printf("client %d: client read and unpack pkt error: %s.", idx, err)
			//break
			continue
		}

		switch p := i.(type) {

		case *network.CmdActiveTestReqPkt:
			log.Printf("client %d: receive a network active request: %v.", idx, p)
			rsp := &network.CmdActiveTestRspPkt{
				Reserved: 1,
			}
			err := c.SendRspPkt(rsp, p.SeqId)
			if err != nil {
				log.Printf("client %d: send network active response error: %s.", idx, err)
				break
			}else{
				log.Printf("client %d: send network active response success: %s.", idx, err)
			}
		case *network.CmdActiveTestRspPkt:
			log.Printf("client %d: receive a network activetest response: %v.", idx, p)

		case *network.CmdTerminateReqPkt:
			log.Printf("client %d: receive a network terminate request: %v.", idx, p)
			rsp := &network.CmdTerminateRspPkt{}
			err := c.SendRspPkt(rsp, p.SeqId)
			if err != nil {
				log.Printf("client %d: send network terminate response error: %s.", idx, err)
				break
			}
		case *network.CmdTerminateRspPkt:
			log.Printf("client %d: receive a network terminate response: %v.", p)
		case *network.CmdUmbrellaInRspPkt:
			log.Printf("client %d: receive a network umbrella in response: %v.", p)
		}
	}
}

var wg sync.WaitGroup

func main() {
	log.Println("Client example start!")
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
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go startAClient(i + 1, sn[i])
	}
	wg.Wait()
	log.Println("Client example ends!")
}
