/*
 * @Author: qiuling
 * @Date: 2019-06-18 15:01:17
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-24 19:17:28
 */
package middleware

import (
	. "artifact/pkg"
	"artifact/pkg/log"
	"net/http"
	"regexp"
	"strings"

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

func middlewareData(userInfo map[string]string, token string, permission int64) Middleware {
	// hostname, _ := os.Hostname()
	// R(userInfo, "userInfo")

	userID, _ := String2Int64(userInfo["user_id"])

	p := userInfo["pid"]
	pid, _ := String2Int64(p)

	info := UserInfo{
		UserID:   userID,
		Code:     userInfo["code"],
		Phone:    userInfo["phone"],
		Jwt:      userInfo["jwt"],
		Avator:   userInfo["avator"],
		Nickname: userInfo["nickname"],
		Pid:      pid,
	}

	middleware := Middleware{
		Permission: permission,
		Token:      token,
		UserID:     userID,
		UserInfo:   info,
	}

	return middleware
}

// matchRoute 匹配路由, 返回匹配到的路由规则
/* func matchRoute(key1 string, key2 string, method1 string, method2 string) string {
	// $pattern = "@^" . preg_replace('/\\\:[a-zA-Z0-9\_\-]+/', '([a-zA-Z0-9\-\_]+)', preg_quote($key2)) . "$@D";
	// if ($method1 == $method2 && preg_match($pattern, $key1, $matches)) {
	// 	return $key2;
	// }``

	// return "";
	// pattern := "@^" + regexp.ReplaceAllStringFunc('/\\\:[a-zA-Z0-9\_\-]+/', '([a-zA-Z0-9\-\_]+)')

	R(key1, "key1")

	re, _ := regexp.Compile("/:[0-9a-zA-Z_]+")

	rep2 := re.ReplaceAllString(key2+"$@D", "([a-zA-Z0-9\_\-]+)")

	R(rep2, "rep2")

	return key2
} */

// KeyMatch determines whether key1 matches the pattern of key2 (similar to RESTful path), key2 can contain a *.
// For example, "/foo/bar" matches "/foo/*", "/resource1" matches "/:resource"
func KeyMatch(key1 string, key2 string) bool {
	key2 = strings.Replace(key2, "/*", "/.*", -1)

	re := regexp.MustCompile(`(.*):[^/]+(.*)`)
	for {
		if !strings.Contains(key2, "/:") {
			break
		}

		key2 = re.ReplaceAllString(key2, "$1[^/]+$2")
	}

	return RegexMatch(key1, "^"+key2+"$")
}

// RegexMatch determines whether key1 matches the pattern of key2 in regular expression.
func RegexMatch(key1 string, key2 string) bool {
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		log.Warn(err)
		// panic(err)
	}
	return res
}
