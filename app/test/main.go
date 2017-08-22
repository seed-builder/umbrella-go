package main

import (
	//"umbrella/network"
	"fmt"
	//"strconv"
	"encoding/hex"
)
func main()  {
	//conn := &network.Conn{}
	//data := []byte{ 0xAA, 0x04, 0xAA, 0x81, 0x2F, 0x55, 0xAA, 0x04, 0xAA, 0x88, 0x4F, 0x88, 0x4F,0x88, 0x4F,0x88, 0x4F,0x88, 0x4F, 0x55}
	//result := conn.ParsePkt(20, data)
	//fmt.Printf("parse pkt is : %x ", result)
	//v := "8804e384"
	//x := fmt.Sprintf("%X", "8804e384")
	//fmt.Println(x)
	//
	//var dst []byte
	//fmt.Sscanf(x, "%X", &dst)
	//fmt.Println(string(dst))

	x := fmt.Sprintf("%X", []byte{0x88, 0x04, 0xe3, 0x84})
	fmt.Println(x)
    buf,_ := hex.DecodeString("8804e384")
	fmt.Printf("%x", buf)
}