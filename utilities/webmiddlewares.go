package utilities

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
)

//cors 跨域
func Cores() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

//verify 签名验证
func VerifySign() gin.HandlerFunc {
	return func(c *gin.Context) {
		sn :=  c.Param("sn")
		customerId := c.Param("customerId")
		sign := c.Query("sign")
		psd := fmt.Sprintf("%s%s",  customerId, sn)
		sign2 := Md5Encrypt([]byte(psd), []byte(SysConfig.Salt))
		if sign == sign2 {
			c.Next()
		}else{
			c.JSON(http.StatusOK, gin.H{"success": false, "err": "sign error, s = 【"+psd+"】, the sign is 【"+ sign2 +"】 " })
			c.Abort()
		}
	}
}

