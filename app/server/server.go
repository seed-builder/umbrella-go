package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
	"umbrella"
	"umbrella/utilities"
	"umbrella/network"
	"os"
	"bufio"
	"strings"
)

func restApi()  {
	utilities.SysLog.Info("REST api 服务启动, 端口：", utilities.SysConfig.HttpPort)
	r := gin.Default()
	r.Use(utilities.Cores())
	// This handler will match /user/john but will not match neither /user/ or /user
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!" )
	})
	umbrella.LoadEquipmentRoutes(r)

	s := &http.Server{
		Addr:           ":" + utilities.SysConfig.HttpPort,
		Handler:        r,
		ReadTimeout:    90 * time.Second,
		WriteTimeout:   90 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}


func equipmentSrv(){
	//defer wg.Done()

	addr := utilities.SysConfig.TcpIp + ":" + utilities.SysConfig.TcpPort
	duration := time.Duration(utilities.SysConfig.TcpTestTimeout)* time.Second
	utilities.SysLog.Info("设备监听服务启动, 地址： ", addr)
	err := umbrella.EquipmentSrv.ListenAndServe(
			addr,
			network.V10,
			duration,
			utilities.SysConfig.TcpTestMax,
			nil,
		)
	if err != nil {
		utilities.SysLog.Info("设备监听服务启动失败：", err)
	}
}

func main() {
	utilities.SysLog.Info("服务开始启动....")
	go restApi()
	go equipmentSrv()

	utilities.SysLog.Info("请输入命令...")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if arr := strings.Fields(line); len(arr) > 0 {
			switch arr[0] {
			case "exit":
				utilities.SysLog.Info("系统退出")
				umbrella.EquipmentSrv.Close()
				os.Exit(0)
			case "open":
				umbrella.EquipmentSrv.OpenChannel(arr[1], 0)
			default:
				utilities.SysLog.Info("请输入正确的命令")
			}
		}
		utilities.SysLog.Info("请输入命令...")
	}
	if err := scanner.Err(); err != nil {
		utilities.SysLog.Error("命令错误:", err)
	}
}
