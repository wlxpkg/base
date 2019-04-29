package model

import (
	"bytes"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gookit/config"
	"github.com/jinzhu/gorm"
)

var Eloquent *gorm.DB

func init() {
	var err error
	mysqlLink := bytes.NewBufferString("")

	mysqlLink.WriteString(config.GetEnv("MYSQL_USERNAME", "root"))
	mysqlLink.WriteString(":" + config.GetEnv("MYSQL_PASSWORD", "root") + "@tcp")
	mysqlLink.WriteString("(" + config.GetEnv("MYSQL_HOST", "127.0.0.1"))
	mysqlLink.WriteString(":" + config.GetEnv("MYSQL_PORT", "3306") + ")")
	mysqlLink.WriteString("/" + config.GetEnv("MYSQL_DATABASE", "artifact"))
	mysqlLink.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=10ms")

	fmt.Printf("%T\n", mysqlLink.String()) // true

	Eloquent, err = gorm.Open("mysql", mysqlLink.String())
	if err != nil {
		fmt.Printf("\nmysql connect err %v\n", err)
	}

	if Eloquent.Error != nil {
		fmt.Printf("\ndatabase err %v\n", Eloquent.Error)
	}
	Eloquent.DB().SetMaxIdleConns(10)
	//Eloquent.DB().SetMaxOpenConns(1000)
}
