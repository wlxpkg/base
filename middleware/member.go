/*
 * @Author: qiuling
 * @Date: 2019-06-20 16:58:11
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-19 15:14:10
 */
package middleware

import (
	. "artifact/pkg"
	"artifact/pkg/model"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func Member() gin.HandlerFunc {
	return func(c *gin.Context) {

		userInfo, err := getUser(c)
		userID, _ := String2Int64(userInfo["user_id"])

		if err != nil {
			err = errors.New("ERR_INVALID_TOKEN")
			Abort(c, err)
			return
		}

		if userID == 0 {
			c.Set("middleware", &Middleware{})
			c.Next()
		}

		permission := getPermission(c, userID)

		middleware := middlewareData(userInfo, permission)
		// R(middleware, "middleware")

		// 设置 example 变量
		c.Set("middleware", middleware)
		c.Next()
	}
}

// getPermission 检测会员是否有权限
// <0 有权限无限次数, >0 有权限有限次数, =0 无权限
func getPermission(c *gin.Context, userID int64) int64 {
	if userID == 0 {
		return 0
	}
	path := c.Request.URL.Path
	method := c.Request.Method

	route := getRoute(path, method)
	if route == "" {
		// 无需鉴权则直接返回 0
		return 0
	}

	clientID := c.GetHeader("client-id")

	permission := checkRole(userID, route, method, clientID)
	return permission
}

// getRoute 获取本次请求匹配的路由
func getRoute(path string, method string) (route string) {
	allRoute := model.MemberRoute()

	route = ""

	for _, routes := range allRoute {
		routeSli := strings.Split(routes, "@")
		// R(routeSli, "routeSli")
		// route := matchRoute(path, routeSli[1], method, routeSli[0])
		if KeyMatch(path, routeSli[1]) && method == routeSli[0] {
			route = routeSli[1]
			return
		}
	}

	return
}

func checkRole(userID int64, route string, method string, clientID string) (count int64) {
	// return 0
	roleIds := model.GetRoleIds(route, method)
	count = 0

	for _, roleId := range roleIds {
		memberRole := model.MemberRole(userID, roleId, clientID)

		if memberRole == -1 {
			count = -1
			return
		} else if memberRole > 0 {
			count += memberRole
		} else {
			continue
		}
	}
	return
}
