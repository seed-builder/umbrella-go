package umbrella

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func LoadEquipmentRoutes(r gin.IRouter)  {
	r.GET("/equipment", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})
	r.POST("/equipment/:sn/open-channel", func(c *gin.Context) {
		sn :=  c.Param("sn")
		success, err := EquipmentSrv.OpenChannel(sn)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"success": success, "err": "" })
		}else{
			c.JSON(http.StatusOK, gin.H{"success": success, "err": err.Error() })
		}

	})
}

func LoadUmbrellaRoutes(r gin.IRouter)  {
	r.GET("/umbrella", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})

}