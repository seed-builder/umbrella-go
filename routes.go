package umbrella

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"umbrella/models"
	"strconv"
	"strings"
	"fmt"
	"umbrella/utilities"
	"umbrella/network"
)

func LoadEquipmentRoutes(r gin.IRouter)  {

	r.GET("/equipment", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})

	r.POST("/equipment/:sn/opennum/:num", func(c *gin.Context) {
		sn :=  c.Param("sn")
		num :=  c.Param("num")
		cn, er := strconv.Atoi(num)
		if er != nil{
			cn = 0
		}
		conn , ok:= EquipmentSrv.EquipmentConns[sn]
		if !ok {
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "设备离线" })
			return
		}
		channelNum, seqId, err := EquipmentSrv.OpenChannel(sn, uint8(cn))
		if err == nil {
			chan_sn, ok := conn.UmbrellaRequests[seqId]
			if ok {
				umbrellaRequest := <-chan_sn
				if !umbrellaRequest.Success {
					c.JSON(http.StatusOK, gin.H{"success": false, "err": umbrellaRequest.Err})
					return
				}
				conn.Equipment.OutChannel(channelNum, nil)
				umbrella := &models.Umbrella{}
				//sn := strconv.Itoa(int(umbrellaSn))
				umbrellaSn := strings.ToUpper(umbrellaRequest.Sn)
				umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)

				c.JSON(http.StatusOK, gin.H{"success": true, "equipment_sn": sn, "channel_num": channelNum, "umbrella_sn": umbrellaSn, "err": ""})
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
	})

	r.POST("/equipment/:sn/set-channel/:num", func(c *gin.Context) {
		sn :=  c.Param("sn")
		num :=  c.Param("num")
		validStr := c.PostForm("valid")
		cn, er := strconv.Atoi(num)
		valid, er := strconv.ParseBool(validStr)
		if er != nil{
			cn = 0
		}
		_ , ok:= EquipmentSrv.EquipmentConns[sn]
		if !ok {
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "设备离线" })
			return
		}
		EquipmentSrv.SetChannel(sn, uint8(cn), valid)

		c.JSON(http.StatusOK, gin.H{"success": true, "err": nil })
	})

	r.POST("/customer/:customerId/hire/:sn", utilities.VerifySign(), func(c *gin.Context) { //
		sn :=  c.Param("sn")
		customerId := c.Param("customerId")
		sign := c.Query("sign")
		fmt.Println("sign = ", sign)
		cid, _ := strconv.ParseUint(customerId, 10, 32)

		utilities.SysLog.Infof("收到客户【%s】在设备【%s】的借伞请求", customerId, sn)

		conn , ok:= EquipmentSrv.EquipmentConns[sn]
		if !ok || conn.State == network.CONN_CLOSED {
			utilities.SysLog.Infof("客户【%s】在设备【%s】离线,无法完成借伞请求", customerId, sn)
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "设备离线" })
			return
		}
		//
		customer := &models.Customer{}
		res := customer.CanBorrowFromEquipment(uint(cid), conn.Equipment.PriceId)
		if !res {
			utilities.SysLog.Infof("客户【%s】押金不足,无法完成借伞请求", customerId)
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "押金不足" })
			return
		}
		channelNum, seqId, err := EquipmentSrv.BorrowUmbrella(uint(cid), sn, 0)
		if err != nil {
			utilities.SysLog.Infof("客户【%s】在设备【%s】借伞请求错误：【%s】", customerId, sn, err.Error())
			c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
			return
		}
		if channelNum == 0 && seqId == 0 {
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "无可用通道" })
			return
		}
		if err == nil {
			chan_sn, ok := conn.UmbrellaRequests[seqId]
			if ok {
				umbrellaRequest := <- chan_sn
				//设置设备当前状态为：等待
				defer func() {
					conn.RunStatus = network.RUN_STATUS_WAITING
				}()

				if !umbrellaRequest.Success {
					utilities.SysLog.Infof("客户【%s】在设备【%s】借伞请求错误：【%s】", customerId, sn, umbrellaRequest.Err)
					c.JSON(http.StatusOK, gin.H{"success": false, "err": umbrellaRequest.Err })
					return
				}
				if conn, ok := EquipmentSrv.EquipmentConns[sn]; ok {
					tx := utilities.MyDB.Begin()
					conn.Equipment.OutChannel(channelNum, tx)
					umbrella := &models.Umbrella{}
					umbrella.InitDb(tx)
					//umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)
					//usn := umbrellaSn
					umbrella.OutEquipment(conn.Equipment, umbrellaRequest.Sn, channelNum)

					hire, err := models.CreateCustomerHire(conn.Equipment, umbrella, uint(cid), tx)
					if hire != nil {
						//hire.FreezeDepositFee(umbrella)
						tx.Commit()
						utilities.SysLog.Infof("客户【%s】在设备【%s】借伞请求成功", customerId, sn)
						c.JSON(http.StatusOK, gin.H{"success": true, "hire_id": hire.ID, "channel": channelNum, "err": ""})
						return
					}else{
						tx.Rollback()
						utilities.SysLog.Infof("客户【%s】在设备【%s】借伞请求错误：【%s】", customerId, sn, err.Error())
						c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
					}
				}
			}
		}
	})

	r.GET("/monitor", func(c *gin.Context) {
		data := gin.H{}
		for sn, conn := range EquipmentSrv.EquipmentConns{
			data[sn] = conn.Equipment
		}
		c.JSON(http.StatusOK, data)
	})
}
