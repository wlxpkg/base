/*
 * @Author: qiuling
 * @Date: 2019-06-17 19:32:28
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-30 10:06:25
 */

package biz

import (
	redis "artifact/pkg/cache"
)

var cache = redis.NewCache().SetPrefix("user")

/* type User struct {
	UserId   string `json:"user_id"`
	Phone    string `json:"phone"`
	Code     string `json:"code"`
	Realname string `json:"realname"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Pid      string `json:"pid"`
} */

func TokenGetUser(uid string) (userInfo map[string]string) {
	baseKey := "base:" + uid
	userInfo = cache.HGetAll(baseKey)
	return
}
