package apiroutes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadEquipmentRoutes(r gin.IRouter)  {
	r.GET("/equipment", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})
}