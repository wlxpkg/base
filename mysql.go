package pkg

import (
	. "artifact/pkg/config"
	"artifact/pkg/log"
	"bytes"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func init() {
	DB = newDB()
}

func newDB() (orm *gorm.DB) {
	// var orm *gorm.DB
	var err error

	mysqlLink := bytes.NewBufferString("")

	mysqlLink.WriteString(Config.Mysql.Username)
	mysqlLink.WriteString(":" + Config.Mysql.Password + "@tcp")
	mysqlLink.WriteString("(" + Config.Mysql.Host)
	mysqlLink.WriteString(":" + Config.Mysql.Port + ")")
	mysqlLink.WriteString("/" + Config.Mysql.Database)
	mysqlLink.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=100ms")

	for orm, err = gorm.Open("mysql", mysqlLink.String()); err != nil; {
		log.Err("mysql connect err: " + err.Error())
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open("mysql", mysqlLink.String())
	}

	if orm.Error != nil {
		log.Err("database err: " + orm.Error.Error())
	}
	// 全局禁用表名复数 TableName不受影响
	orm.SingularTable(true)
	orm.DB().SetMaxIdleConns(100)
	//orm.DB().SetMaxOpenConns(1000)

	return
}
