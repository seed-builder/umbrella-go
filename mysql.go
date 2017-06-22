package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"umbrella/utl"
)

var MyDB *gorm.DB

func init() {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		utl.SysConfig.DbUser,
		utl.SysConfig.DbPassword,
		utl.SysConfig.DbServer,
		utl.SysConfig.DbDatabase,
	)

	MyDB, err = gorm.Open("mysql", dsn)
	// 连接池
	if err == nil {
		MyDB.DB().SetMaxIdleConns(50)
		MyDB.DB().SetMaxOpenConns(100)
		MyDB.DB().Ping()
		MyDB.LogMode(true)
	} else {
		log.Panic("Gorm Open Error: ", err)
	}
}