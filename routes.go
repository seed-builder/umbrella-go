package umbrella

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"umbrella/models"
	"strconv"
)

func LoadEquipmentRoutes(r gin.IRouter)  {

	r.GET("/equipment", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})

	r.POST("/equipment/:sn/open", func(c *gin.Context) {
		sn :=  c.Param("sn")
		channelNum, seqId, err := EquipmentSrv.OpenChannel(sn)
		if err == nil {
			chan_sn, ok := EquipmentSrv.Requests[seqId]
			if ok {
				umbrellaSn := <- chan_sn
				log.Println("channel opened!")

				if conn, ok := EquipmentSrv.EquipmentConns[sn]; ok {
					conn.Equipment.OutChannel(channelNum)
					umbrella := &models.Umbrella{}
					sn := strconv.Itoa(int(umbrellaSn))
					umbrella.OutEquipment(conn.Equipment, sn, channelNum)
				}
				c.JSON(http.StatusOK, gin.H{"success": true, "equipment_sn": sn, "channel_num": channelNum, "umbrella_sn": umbrellaSn,  "err": "" })
				return
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
	})

	r.POST("/customer/:customerId/hire/:sn", func(c *gin.Context) {
		sn :=  c.Param("sn")
		customerId := c.Param("customerId")
		channelNum, seqId, err := EquipmentSrv.OpenChannel(sn)
		if err == nil {
			chan_sn, ok := EquipmentSrv.Requests[seqId]
			if ok {
				umbrellaSn := <- chan_sn
				log.Println("channel opened!")

				if conn, ok := EquipmentSrv.EquipmentConns[sn]; ok {
					conn.Equipment.OutChannel(channelNum)
					umbrella := &models.Umbrella{}
					//umbrella.OutEquipment(conn.Equipment, umbrellaSn, channelNum)
					usn := strconv.Itoa(int(umbrellaSn))
					umbrella.OutEquipment(conn.Equipment, usn, channelNum)
					hire := &models.CustomerHire{}
					cid, _ := strconv.ParseUint(customerId, 10, 32)
					hire.Create(conn.Equipment, umbrella, uint(cid))

					c.JSON(http.StatusOK, gin.H{"success": true, "hire_id": hire.ID , "err": "" })
					return
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": false, "err": err.Error() })
	})
}
