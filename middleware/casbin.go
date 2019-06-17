/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:33:04
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-17 19:10:16
 */

package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Casbin() gin.HandlerFunc {
	return func(c *gin.Context) {

		uid := getUser(c)

		fmt.Printf("%+v\n", uid)

		// 设置 example 变量
		c.Set("example", uid)

		c.Next()
	}
}

func getUser(c *gin.Context) (err error) {
	authorization := c.GetHeader("authorization")
	jwt := strings.TrimPrefix(authorization, "Bearer ")

	fmt.Printf("%+v\n", jwt)

	if jwt == "" {
		err = errors.New("ERR_INVALID_TOKEN")
	}
	return
	// method := r.Method
	// path := r.URL.Path
}
