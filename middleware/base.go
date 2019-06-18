/*
 * @Author: qiuling
 * @Date: 2019-06-18 15:01:17
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-18 15:46:45
 */
package middleware

import (
	. "artifact/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Abort(c *gin.Context, e error) {

	errors, ok := Errs[e.Error()]
	if !ok {
		errors = Errs["ERR_UNKNOW_ERROR"]
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    errors.Code,
		"message": errors.Message,
		"data":    "",
	})
	c.Abort()
}

func GetPolicy() {

}
