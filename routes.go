package umbrella

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"umbrella/models"
	"strconv"
	"strings"
	"fmt"
	"umbrella/utilities"
)

func LoadEquipmentRoutes(r gin.IRouter)  {

	r.GET("/equipment", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})

	r.POST("/equipment/:sn/open", func(c *gin.Context) {
		sn :=  c.Param("sn")
		channelNum, seqId, err := EquipmentSrv.OpenChannel(sn, 0)
		if err == nil {
			chan_sn, ok := EquipmentSrv.Requests[seqId]
			if ok {
				umbrellaSn := <- chan_sn
				log.Println("channel opened!")

				if conn, ok := EquipmentSrv.EquipmentConns[sn]; ok {
					conn.Equipment.OutChannel(channelNum, nil)
					umbrella := &models.Umbrella{}
					//sn := strconv.Itoa(int(umbrellaSn))
					umbrellaSn = strings.ToUpper(umbrellaSn)
					umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)
				}
				c.JSON(http.StatusOK, gin.H{"success": true, "equipment_sn": sn, "channel_num": channelNum, "umbrella_sn": umbrellaSn,  "err": "" })
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
	})

	r.POST("/equipment/:sn/opennum/:num", func(c *gin.Context) {
		sn :=  c.Param("sn")
		num :=  c.Param("num")
		cn, er := strconv.Atoi(num)
		if er != nil{
			cn = 0
		}
		channelNum, seqId, err := EquipmentSrv.OpenChannel(sn, uint8(cn))
		if err == nil {
			chan_sn, ok := EquipmentSrv.Requests[seqId]
			if ok {
				umbrellaSn := <- chan_sn
				log.Println("channel opened! umbrellaSn : ", umbrellaSn)
				if umbrellaSn == ""{
					c.JSON(http.StatusOK, gin.H{"success": false, "err": "超时" })
					return
				}
				if conn, ok := EquipmentSrv.EquipmentConns[sn]; ok {
					conn.Equipment.OutChannel(channelNum, nil)
					umbrella := &models.Umbrella{}
					//sn := strconv.Itoa(int(umbrellaSn))
					umbrellaSn = strings.ToUpper(umbrellaSn)
					umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)
				}
				c.JSON(http.StatusOK, gin.H{"success": true, "equipment_sn": sn, "channel_num": channelNum, "umbrella_sn": umbrellaSn,  "err": "" })
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
	})

	r.POST("/customer/:customerId/hire/:sn", func(c *gin.Context) { //utilities.VerifySign(),
		sn :=  c.Param("sn")
		customerId := c.Param("customerId")
		sign := c.Query("sign")
		fmt.Println("sign = ", sign)
		channelNum, seqId, err := EquipmentSrv.OpenChannel(sn, 0)
		if channelNum == 0 && seqId == 0 {
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "无可用通道" })
			return
		}
		if err == nil {
			chan_sn, ok := EquipmentSrv.Requests[seqId]
			if ok {
				umbrellaSn := <- chan_sn
				if umbrellaSn == ""{
					c.JSON(http.StatusOK, gin.H{"success": false, "err": "超时" })
					return
				}
				if conn, ok := EquipmentSrv.EquipmentConns[sn]; ok {
					tx := utilities.MyDB.Begin()
					conn.Equipment.OutChannel(channelNum, tx)
					umbrella := &models.Umbrella{}
					umbrella.InitDb(tx)
					//umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)
					//usn := umbrellaSn
					umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)
					cid, _ := strconv.ParseUint(customerId, 10, 32)
					hire, err := models.CreateCustomerHire(conn.Equipment, umbrella, uint(cid), tx)
					if hire != nil {
						//hire.FreezeDepositFee(umbrella)
						tx.Commit()
						c.JSON(http.StatusOK, gin.H{"success": true, "hire_id": hire.ID, "channel": channelNum, "err": ""})
						return
					}else{
						tx.Rollback()
						c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
					}
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
	})

	r.GET("/monitor", func(c *gin.Context) {
		data := gin.H{}
		for sn, conn := range EquipmentSrv.EquipmentConns{
			data[sn] = conn.Equipment
		}
		c.JSON(http.StatusOK, data)
	})
}
