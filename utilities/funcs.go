package utilities

import (
	"math"
	"crypto/md5"
	"encoding/hex"
)

func Round(f float64, prec int) float64 {
	pow10_n := math.Pow10(prec)
	return math.Trunc((f + 0.5/pow10_n) * pow10_n ) / pow10_n
}

func Md5Encrypt(password []byte, salt []byte) string {
	h := md5.New()
	h.Write(password) // 需要加密的字符串为 123456
	if salt != nil {
		h.Write(salt)
	}
	cipherStr := h.Sum(nil)
	//fmt.Println(cipherStr)
	return hex.EncodeToString(cipherStr) // 输出加密结果
}