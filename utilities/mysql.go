package utilities

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
)

var MyDB *gorm.DB

func init() {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		SysConfig.DbUser,
		SysConfig.DbPassword,
		SysConfig.DbServer,
		SysConfig.DbDatabase,
	)

	MyDB, err = gorm.Open("mysql", dsn)
	// 连接池
	if err == nil {
		MyDB.DB().SetMaxIdleConns(50)
		MyDB.DB().SetMaxOpenConns(100)
		MyDB.DB().Ping()
		MyDB.LogMode(true)
	} else {
		SysLog.Panic("Gorm Open Error: ", err)
	}
}