package main

import (
	"umbrella/models"
	"fmt"
)
func main()  {
	fmt.Println("begin ...!")
	ch := &models.CustomerHire{}
	eq := &models.Equipment{}
	eq.Query().Find(eq, "id=?", 1)
	umbrella := &models.Umbrella{}
	umbrella.Query().Find(umbrella, "id=?", 1)
	ch.Create(eq, umbrella, 5)
	fmt.Println("complete", ch)
}