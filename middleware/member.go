/*
 * @Author: qiuling
 * @Date: 2019-06-20 16:58:11
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:10:17
 */
package middleware

import (
	"bytes"
	"errors"
	. "github.com/wlxpkg/base"
	"github.com/wlxpkg/base/model"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Member() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 后置数据准备
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()                                        //  must close
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新赋值

		c.Set("bodyCopy", string(bodyBytes))

		userInfo, err := getUser(c)
		userID, _ := String2Int64(userInfo["user_id"])

		// R(userInfo, "userInfo")
		if err != nil {
			err = errors.New("ERR_INVALID_TOKEN")
			Abort(c, err)
			return
		}

		guestAllow := checkGuest(c, userID)
		if !guestAllow {
			err = errors.New("ERR_UNLOGIN")
			Abort(c, err)
			return
		}

		if userID == 0 || userID == 403 {
			c.Set("middleware", &Middleware{})
			c.Next()
		}

		permission := getPermission(c, userID)

		var flower int64
		if userInfo["has_flower"] == "1" {
			flower = -1
		} else {
			flower = model.GetFlower(userID)
		}

		middleware := middlewareData(userInfo, permission, flower)
		// R(middleware, "middleware")

		// 设置 example 变量
		c.Set("middleware", middleware)
		c.Next()
	}
}

func checkGuest(c *gin.Context, userID int64) bool {
	if userID == 0 || userID == 403 {
		path := c.Request.URL.Path
		method := c.Request.Method

		route := getRoute(path, method, 0)
		if route == "" {
			// 没匹配到白名单, 则无权限, 需要登录
			return false
		}
	}

	return true
}

// getPermission 检测会员是否有权限
func getPermission(c *gin.Context, userID int64) bool {
	if userID == 0 || userID == 403 {
		return false
	}
	path := c.Request.URL.Path
	method := c.Request.Method

	route := getRoute(path, method, 1)
	if route == "" {
		// 无需鉴权则直接返回 false
		return false
	}

	clientID := c.GetHeader("client-id")

	permission := checkRole(userID, route, method, clientID)
	return permission
}

// getRoute 获取本次请求匹配的路由
// rtype 路由类型, 0:guest, 1:member
func getRoute(path string, method string, rtype int) (route string) {

	var allRoute []string
	switch rtype {
	case 0:
		allRoute = model.GuestRoute()
	case 1:
		allRoute = model.MemberRoute()
	}

	route = ""

	for _, routes := range allRoute {
		routeSli := strings.Split(routes, "@")
		if len(routeSli) < 2 {
			continue
		}

		if KeyMatch(path, routeSli[1]) && method == routeSli[0] {
			route = routeSli[1]
			break
		}
	}

	return
}

func checkRole(userID int64, route string, method string, clientID string) (permission bool) {
	// return 0
	roleIds := model.GetRoleIds(route, method)
	permission = false

	for _, roleId := range roleIds {
		expireAt := model.MemberRole(userID, roleId, clientID)

		if expireAt == "" {
			return false
		}

		if String2Unix(expireAt) > time.Now().Unix() {
			return true
		}
	}
	return
}
