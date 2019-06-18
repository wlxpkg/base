/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:33:04
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-18 15:47:11
 */

package middleware

import (
	"artifact/pkg/biz"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {

		token, userInfo, err := getUser(c)

		if err != nil {
			err = errors.New("ERR_INVALID_TOKEN")
			Abort(c, err)
		}

		e := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")

		fmt.Printf("%+v\n", e)
		// fmt.Printf("%+v\n", userInfo)

		hostname, _ := os.Hostname()
		middleware := map[string]string{
			"token":    token,
			"userID":   userInfo["user_id"],
			"hostname": hostname,
		}

		// 设置 example 变量
		c.Set("middleware", middleware)
		c.Set("userInfo", userInfo)

		c.Next()
	}
}

func getUser(c *gin.Context) (token string, userInfo map[string]string, err error) {
	authorization := c.GetHeader("authorization")
	jwt := strings.TrimPrefix(authorization, "Bearer ")

	fmt.Printf("jwt:%+v\n", jwt)

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
	// method := r.Method
	// path := r.URL.Path
}
