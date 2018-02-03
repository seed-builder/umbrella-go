package test

import (
	"testing"
	"regexp"
	"fmt"
)

func Test_Regex(t *testing.T){
	reg := regexp.MustCompile(`设备【(\w+)】`)
	res := reg.FindSubmatch([]byte("设备【E198402mqvw】通道【1】异常/超时"))
	fmt.Printf("%q len: %d, 1: %q, 2:%q", res, len(res), res[0], res[1])
}