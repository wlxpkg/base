/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:33:04
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-11-09 11:09:24
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
	"time"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/joncalhoun/qson"
)

var e *casbin.Enforcer
var timerTicker *time.Ticker

func init() {
	text :=
		`
	[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act, eft

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow)) && !some(where (p.eft == deny))

	[matchers]
	m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
	`
	m := casbin.NewModel(text)

	a := model.NewAdapter("mysql", mysqlLink(), true)
	e := casbin.NewEnforcer(m, a)
	_ = e.LoadPolicy()

	timeRefresh()
	// fmt.Printf("LoadPolicy ERR: %+v\n", err)
}

func Casbin() gin.HandlerFunc {

	return func(c *gin.Context) {

		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()                                        //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新赋值

		userInfo, err := getUser(c)
		userID := userInfo["user_id"]

		if err != nil {
			err = errors.New("ERR_INVALID_TOKEN")
			Abort(c, err)
			return
		}

		method := c.Request.Method
		path := c.Request.URL.Path

		if ok, _ := e.Enforce(userID, path, method); !ok {
			err = errors.New("ERR_UNAUTHORIZED")
			Abort(c, err)
			return
		}

		middleware := middlewareData(userInfo, true, 0)
		// 设置中间件变量
		c.Set("middleware", middleware)

		// 后置数据准备
		c.Set("adminid", userID)
		c.Set("bodyCopy", string(bodyBytes))
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

	DB.Create(&log)
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

// timeRefresh
// 数据更新在php的服务中, 因此无法使用 watcher, 只能 定时更新 casbin数据
func timeRefresh() {
	timerTicker = time.NewTicker(10 * time.Minute) // 10分钟 定时器
	defer timerTicker.Stop()

	go func() {
		for range timerTicker.C {
			_ = e.LoadPolicy()
			R("", "casbin load data")
		}
	}()
}
