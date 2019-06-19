/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:33:04
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-19 14:59:58
 */

package middleware

import (
	"artifact/pkg/biz"
	. "artifact/pkg/config"
	"artifact/pkg/model"
	"bytes"
	"errors"
	"os"
	"strings"

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

		// fmt.Printf("userID:%+v\n", userID)
		// fmt.Printf("method:%+v\n", method)
		// fmt.Printf("path:%+v\n", path)

		if !e.Enforce(userID, path, method) {
			err = errors.New("ERR_UNAUTHORIZED")
			Abort(c, err)
			return
		}

		hostname, _ := os.Hostname()
		middleware := map[string]string{
			"token":    token,
			"userID":   userID,
			"hostname": hostname,
		}

		// 设置 example 变量
		c.Set("middleware", middleware)
		c.Set("middlewareUserInfo", userInfo)

		c.Next()
	}
}

func getUser(c *gin.Context) (token string, userInfo map[string]string, err error) {
	authorization := c.GetHeader("authorization")
	jwt := strings.TrimPrefix(authorization, "Bearer ")

	// fmt.Printf("jwt:%+v\n", jwt)

	if jwt == "" {
		err = errors.New("ERR_INVALID_TOKEN")
		return
	}

	token, err = biz.Jwt2Token(jwt)

	if token == "" || err != nil {
		err = errors.New("ERR_INVALID_TOKEN")
		return
	}

	userInfo = biz.TokenGetUser(token)
	// fmt.Printf("userInfo:%+v\n", userInfo)

	return
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
