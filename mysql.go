package pkg

import (
	. "artifact/pkg/config"
	"bytes"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	mysqlLink := bytes.NewBufferString("")

	mysqlLink.WriteString(Config.Mysql.Username)
	mysqlLink.WriteString(":" + Config.Mysql.Password + "@tcp")
	mysqlLink.WriteString("(" + Config.Mysql.Host)
	mysqlLink.WriteString(":" + Config.Mysql.Port + ")")
	mysqlLink.WriteString("/" + Config.Mysql.Database)
	mysqlLink.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=100ms")

	DB, err = gorm.Open("mysql", mysqlLink.String())
	if err != nil {
		fmt.Printf("\nmysql connect err %v\n", err)
	}

	if DB.Error != nil {
		fmt.Printf("\ndatabase err %v\n", DB.Error)
	}
	// set table prefix
	// gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
	// 	return "course_" + defaultTableName
	// }

	// 全局禁用表名复数 TableName不受影响
	DB.SingularTable(true)

	DB.DB().SetMaxIdleConns(100)
	//DB.DB().SetMaxOpenConns(1000)
}
