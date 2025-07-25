package base

import (
	"bytes"
	"time"

	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DB gorm.DB
var DB *gorm.DB

// ZwydDB gorm.DB for zwyd
var ZwydDB *gorm.DB

func init() {
	DB = newDB()

	if Config.Mysql.Dump == "true" {
		DB.LogMode(true)
	}
	ZwydDB = newZwydDB()
	if Config.Mysql.Dump == "true" {
		ZwydDB.LogMode(true)
	}
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
	mysqlLink.WriteString("?charset=utf8mb4&parseTime=True&loc=Local&timeout=100ms")

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
	orm.DB().SetConnMaxLifetime(55 * time.Second)
	// orm.DB().SetMaxOpenConns(1000)

	return
}

func newZwydDB() (orm *gorm.DB) {
	// var orm *gorm.DB
	var err error

	mysqlLink := bytes.NewBufferString("")

	mysqlLink.WriteString(Config.MysqlZWYD.Username)
	mysqlLink.WriteString(":" + Config.MysqlZWYD.Password + "@tcp")
	mysqlLink.WriteString("(" + Config.MysqlZWYD.Host)
	mysqlLink.WriteString(":" + Config.MysqlZWYD.Port + ")")
	mysqlLink.WriteString("/" + Config.MysqlZWYD.Database)
	mysqlLink.WriteString("?charset=utf8mb4&parseTime=True&loc=Local&timeout=100ms")

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
	orm.DB().SetConnMaxLifetime(55 * time.Second)
	// orm.DB().SetMaxOpenConns(1000)

	return
}
