package pkg

import (
	"bytes"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/config"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	mysqlLink := bytes.NewBufferString("")

	mysqlLink.WriteString(config.GetEnv("MYSQL_USERNAME", "root"))
	mysqlLink.WriteString(":" + config.GetEnv("MYSQL_PASSWORD", "root") + "@tcp")
	mysqlLink.WriteString("(" + config.GetEnv("MYSQL_HOST", "127.0.0.1"))
	mysqlLink.WriteString(":" + config.GetEnv("MYSQL_PORT", "3306") + ")")
	mysqlLink.WriteString("/" + config.GetEnv("MYSQL_DATABASE", "artifact"))
	mysqlLink.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=100ms")

	fmt.Printf("%T\n", mysqlLink.String()) // true

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
