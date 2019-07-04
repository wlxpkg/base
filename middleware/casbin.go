/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:33:04
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-04 19:05:44
 */

package middleware

import (
	// . "artifact/pkg"
	. "artifact/pkg/config"
	"artifact/pkg/model"
	"bytes"
	"errors"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func Casbin() gin.HandlerFunc {

	a := model.NewAdapter("mysql", mysqlLink(), true)
	e := casbin.NewEnforcer("../vendor/artifact/pkg/middleware/lauthz-rbac-model.conf", a)
	_ = e.LoadPolicy()
	// fmt.Printf("LoadPolicy ERR: %+v\n", err)

	return func(c *gin.Context) {

		token, userInfo, err := getUser(c)
		userID := userInfo["user_id"]

		if err != nil {
			err = errors.New("ERR_INVALID_TOKEN")
			Abort(c, err)
			return
		}

		method := c.Request.Method
		path := c.Request.URL.Path

		if !e.Enforce(userID, path, method) {
			err = errors.New("ERR_UNAUTHORIZED")
			Abort(c, err)
			return
		}

		middleware := middlewareData(userInfo, token, 1)
		// 设置 example 变量
		c.Set("middleware", middleware)

		c.Next()
	}
}

func mysqlLink() string {
	mysqlLink := bytes.NewBufferString("")

	mysqlLink.WriteString(Config.Mysql.Username)
	mysqlLink.WriteString(":" + Config.Mysql.Password + "@tcp")
	mysqlLink.WriteString("(" + Config.Mysql.Host)
	mysqlLink.WriteString(":" + Config.Mysql.Port + ")")
	mysqlLink.WriteString("/" + Config.Mysql.Database)
	// mysqlLink.WriteString("?charset=utf8&parseTime=True&loc=Local&timeout=100ms")

	return mysqlLink.String()
}
