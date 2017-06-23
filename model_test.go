package main

import (
	"testing"
	"umbrella/models"
	//"umbrella/utility"
)

func Test_Create(t *testing.T) {
	site := models.Site{ Name: "test22", Province:"test22", City: "test22", District: "tst22"}
	//site.Entity = site
	site.Save()
	t.Log("test end !!")
	t.Log(site)

	//utility.MyDB.Save(&site)
	//t.Log("db save ")
	//t.Log(site)
}
