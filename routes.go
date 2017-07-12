package umbrella

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"
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

func LoadCustomerHireRoutes(r gin.IRouter)  {
	r.POST("/customer-hire/:id/do", func(c *gin.Context) {
		id := c.Param("id")
		hire_id, _ := strconv.Atoi(id)
		success, err := EquipmentSrv.DoHire(uint(hire_id))
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"success": success, "err": "" })
		}else{
			c.JSON(http.StatusOK, gin.H{"success": success, "err": err.Error() })
		}
	})

}