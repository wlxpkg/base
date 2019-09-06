/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:33:04
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-09-06 16:05:13
 */

package middleware

import (
	. "artifact/pkg"
	. "artifact/pkg/config"
	"artifact/pkg/model"
	"bytes"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/joncalhoun/qson"
)

func Casbin() gin.HandlerFunc {

	a := model.NewAdapter("mysql", mysqlLink(), true)
	e := casbin.NewEnforcer("../vendor/artifact/pkg/middleware/lauthz-rbac-model.conf", a)
	_ = e.LoadPolicy()
	// fmt.Printf("LoadPolicy ERR: %+v\n", err)

	return func(c *gin.Context) {

		userInfo, err := getUser(c)
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

		middleware := middlewareData(userInfo, true, 0)
		// 设置中间件变量
		c.Set("middleware", middleware)

		// 后置数据准备
		c.Set("adminid", userID)
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()                                        //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新赋值

		// c.Set("bodyCopy", body)
		// c.Set("dataType", dataType)

		// 执行业务
		c.Next()

		// 后置中间件
		if method != "GET" {
			go addLog(c, userID, string(bodyBytes))
		}
	}
}

func addLog(c *gin.Context, adminId string, bodyString string) {
	path := c.Request.URL.Path
	method := c.Request.Method

	// R(bodyString, "bodyString")
	var bodyBytes []byte
	var bodyData map[string]interface{}

	contentType := c.ContentType()
	// R(contentType, "contentType")
	if strings.Contains(contentType, "form-urlencoded") {
		bodyBytes, _ = qson.ToJSON(bodyString)
	} else if strings.Contains(contentType, "json") {
		// 转
		bodyData, _ = JsonDecode(bodyString)
		bodyBytes, _ = JsonEncode(bodyData)
	}

	if bodyBytes == nil {
		bodyBytes = []byte("[]")
	}

	log := model.AdminOperationLog{
		UserId:  adminId,
		Path:    "/" + Config.Redis.Prefix + path,
		Method:  method,
		Ip:      c.ClientIP(),
		Request: bodyBytes,
	}
	// R(log, "log")

	DB.Debug().Create(&log)
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
