package main

import (
	//"umbrella/network"
	"fmt"
	//"strconv"
	//"encoding/hex"
	"umbrella/network"
)
func main()  {
	conn := &network.Conn{}
	data := []byte{ 0xaa, 0x06, 0xc6, 0x84, 0x04, 0x01, 0x55, 0x55, 0xaa, 0x06, 0xc6, 0x84, 0x04, 0x01, 0x55, 0x55 }
	result := conn.ParsePkt(16, data)
	fmt.Printf("parse pkt is : %x ", result)
	//
	//var dst []byte
	//fmt.Sscanf(x, "%X", &dst)
	//fmt.Println(string(dst))

	//x := fmt.Sprintf("%X", []byte{0x88, 0x04, 0xe3, 0x84})
	//fmt.Println(x)
	//buf,_ := hex.DecodeString("8804e384")
	//fmt.Printf("%x", buf)
}