package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"umbrella"
	"umbrella/utilities"
	"umbrella/network"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"sync"
)

var wg sync.WaitGroup

func main(){
	fmt.Println("please enter ctl+c to terminate")
	shutdown := make(chan struct{})

	//go restApi()
	go equipmentSrv()

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case c := <-shutdown:
			fmt.Println("shutdown system...", c)
			umbrella.EquipmentSrv.Close()
			return
		}
	}()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT)

	s := <-c
	close(shutdown)
	fmt.Println("Got signal:", s)
	wg.Wait()
	fmt.Println("System Quit")
}

func restApi()  {
	defer wg.Done()

	log.Println("rest api service begin...")
	r := gin.Default()
	r.Use(cors())
	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!" )
	})
	umbrella.LoadEquipmentRoutes(r)
	umbrella.LoadUmbrellaRoutes(r)
	umbrella.LoadCustomerHireRoutes(r)
	r.Run(":" + utilities.SysConfig.HttpPort)
	log.Println("restApi return ")
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func equipmentSrv(){
	defer wg.Done()

	addr := ":" + utilities.SysConfig.TcpPort
	duration := time.Duration(utilities.SysConfig.TcpTestTimeout)* time.Second
	err := umbrella.EquipmentSrv.ListenAndServe(
		addr,
		network.V10,
		duration,
		utilities.SysConfig.TcpTestMax,
		nil,
	)
	if err != nil {
		log.Println("equipment ListenAndServ error:", err)
	}
	log.Println("equipmentSrv return ")
	return
}
