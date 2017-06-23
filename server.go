package main

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"net"
	"time"
	"fmt"
	"os"
	"umbrella/apiroutes"
	"umbrella/utilities"
)

func main(){
	log.Println("umbrella service begin...")
	go restApi()
	machineSvr()
}

func restApi()  {
	log.Println("rest api service begin...")
	router := gin.Default()

	// This handler will match /user/john but will not match neither /user/ or /user
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!" )
	})
	apiroutes.LoadEquipmentRoutes(router)
	apiroutes.LoadUmbrellaRoutes(router)

	router.Run(":" + utilities.SysConfig.HttpPort)

}

func machineSvr(){
	log.Println("machine service begin...")

	service := utilities.SysConfig.TcpPort
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		log.Println("machine receive data!")
		daytime := time.Now().String()
		conn.Write([]byte(daytime)) // don't care about return value
		conn.Close()                // we're finished with this client
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}