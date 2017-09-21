package main

import (
	"umbrella"
	"umbrella/models"
	"fmt"
)
func main()  {
	umbrella.EquipmentSrv.Close()
	fmt.Println("begin ...!")
	//um := &models.Umbrella{}
	//eq := &models.Equipment{}
	//eq.Query().Find(eq, 1)
	//status := um.OutEquipment(eq, "17617E62", 1)
	//fmt.Printf("umbrella out equipment status: %s", status)
	//hire := &models.CustomerHire{}
	//hire.Create(eq, um, 66)
	//hire.FreezeDepositFee(um)

	ch := &models.CustomerHire{}
	ch.UmbrellaReturn(1, 1, 1)
	//time.Sleep(1 * time.Minute)
	fmt.Println("end ...!")
}