/*
 * @Author: qiuling
 * @Date: 2019-08-22 18:26:18
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-08-23 09:33:05
 */
package middleware

import (
	. "artifact/pkg"
	"artifact/pkg/log"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Err(r)
				var err error

				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("内部系统错误")
				}

				errors, ok := Errs[err.Error()]
				if !ok {
					errors = Errors{Code: 0, Message: "内部系统错误"}
				}

				c.JSON(http.StatusOK, gin.H{
					"code":    errors.Code,
					"message": errors.Message,
					"data":    "",
				})
			}
		}()
		c.Next()
	}
}
