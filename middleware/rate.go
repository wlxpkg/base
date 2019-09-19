/*
 * 限流中间件
 * @Author: qiuling
 * @Date: 2019-09-18 16:16:34
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-09-19 11:45:33
 */

package middleware

import (
	. "artifact/pkg"
	"artifact/pkg/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func Rate() gin.HandlerFunc {
	return func(c *gin.Context) {

		authorization := c.GetHeader("authorization")
		jwt := strings.TrimPrefix(authorization, "Bearer ")
		path := c.Request.URL.Path

		check := model.RateCheck(jwt, path)
		if !check {
			err := Excp("ERR_TOO_MANY_REQUEST")
			Abort(c, err)
			return
		}
	}
}
