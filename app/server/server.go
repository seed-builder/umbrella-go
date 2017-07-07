package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"umbrella"
	"umbrella/utilities"
	"umbrella/network"
)

func main(){
	log.Println("umbrella service begin...")
	go restApi()
	equipmentSrv()
}

func restApi()  {
	log.Println("rest api service begin...")
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!" )
	})
	umbrella.LoadEquipmentRoutes(router)
	umbrella.LoadUmbrellaRoutes(router)

	router.Run(":" + utilities.SysConfig.HttpPort)

}

func equipmentSrv(){
	var handlers = []network.Handler{
		network.HandlerFunc(umbrella.HandleConnect),
		network.HandlerFunc(umbrella.HandleUmbrellaIn),
		network.HandlerFunc(umbrella.HandleUmbrellaOut),
	}
	addr := ":" + utilities.SysConfig.TcpPort
	err := umbrella.EquipmentSrv.ListenAndServe(
		addr,
		network.V10,
		5*time.Second,
		3,
		nil,
		handlers...,
	)
	if err != nil {
		log.Println("equipment ListenAndServ error:", err)
	}
	return
}
