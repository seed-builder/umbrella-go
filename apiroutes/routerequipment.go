package apiroutes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"umbrella/services"
)

func LoadEquipmentRoutes(r gin.IRouter)  {
	r.GET("/equipment", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})
	r.POST("/equipment/:sn/open-channel", func(c *gin.Context) {
		sn :=  c.Param("sn")
		success, err := services.EquipmentSrv.OpenChannel(sn)
		c.JSON(http.StatusOK, gin.H{"success": success, "err": err })
	})
}