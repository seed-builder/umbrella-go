package main

import (
	"testing"
	"umbrella/models"
	//"umbrella/utility"
	"github.com/gin-gonic/gin"
)

func Test_Create(t *testing.T) {

	site := models.NewSite(gin.H{ "Name": "test44", "Province":"test44", "City": "test44", "District": "test44"}) //models.Site{ Name: "test22", Province:"test22", City: "test22", District: "tst22"}
	//site := models.Site{ Name: "test33", Province:"test33", City: "test33", District: "test33"}
	//site.Entity = site
	//site.Name = "site222"
	//site.Province = "site222"
	//site.City = "site222"
	//site.District = "site222"
	site.Save()
	t.Log("test end !!")
	t.Log(site)

	//utility.MyDB.Save(&site)
	//t.Log("db save ")
	//t.Log(site)
}
