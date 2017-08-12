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
	"bufio"
	"strings"
)

var wg sync.WaitGroup

func main1(){
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
	//defer wg.Done()

	log.Println("rest api service begin...")
	r := gin.Default()
	r.Use(cors())
	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!" )
	})
	umbrella.LoadEquipmentRoutes(r)
	//r.Run(":" + utilities.SysConfig.HttpPort)

	s := &http.Server{
		Addr:           ":" + utilities.SysConfig.HttpPort,
		Handler:        r,
		ReadTimeout:    90 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func equipmentSrv(){
	//defer wg.Done()

	addr := utilities.SysConfig.TcpIp + ":" + utilities.SysConfig.TcpPort
	duration := time.Duration(utilities.SysConfig.TcpTestTimeout)* time.Second
	log.Println("equipmentSrv start serve at addr: ", addr)
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
}

func main() {
	fmt.Println("server start...! ")

	go restApi()
	go equipmentSrv()

	fmt.Println("please input a command ...")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if arr := strings.Fields(line); len(arr) > 0 {
			switch arr[0] {
			case "exit":
				fmt.Println("server exit...")
				umbrella.EquipmentSrv.Close()
				os.Exit(0)
			case "open":
				umbrella.EquipmentSrv.OpenChannel(arr[1])
			default:

			}
		}
		fmt.Println("please input a command ...")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
