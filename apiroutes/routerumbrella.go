package apiroutes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadUmbrellaRoutes(r gin.IRouter)  {
	r.GET("/umbrella", func(c *gin.Context) {
		c.String(http.StatusOK, "equipment rest api !")
	})

	//伞解锁
	r.POST("/umbrella/:sn/open", func(c *gin.Context) {

	})
}
