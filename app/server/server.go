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

	go restApi()
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
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!" )
	})
	umbrella.LoadEquipmentRoutes(router)
	umbrella.LoadUmbrellaRoutes(router)
	router.Run(":" + utilities.SysConfig.HttpPort)
	log.Println("restApi return ")
}

func equipmentSrv(){
	defer wg.Done()

	var handlers = []network.Handler{
		network.HandlerFunc(umbrella.HandleConnect),
		network.HandlerFunc(umbrella.HandleUmbrellaIn),
		network.HandlerFunc(umbrella.HandleUmbrellaOut),
	}
	addr := ":" + utilities.SysConfig.TcpPort
	duration := time.Duration(utilities.SysConfig.TcpTestTimeout)* time.Second
	err := umbrella.EquipmentSrv.ListenAndServe(
		addr,
		network.V10,
		duration,
		utilities.SysConfig.TcpTestMax,
		nil,
		handlers...,
	)
	if err != nil {
		log.Println("equipment ListenAndServ error:", err)
	}
	log.Println("equipmentSrv return ")
	return
}
