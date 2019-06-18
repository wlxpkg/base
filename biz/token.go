/*
 * @Author: qiuling
 * @Date: 2019-06-17 19:32:28
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-18 14:09:12
 */

package biz

import (
	redis "artifact/pkg/cache"
)

var cache = redis.NewCache().SetPrefix("user")

type User struct {
	UserId         string `json:"user_id"`
	Phone          string `json:"phone"`
	InvitationCode string `json:"invitation_code"`
	Realname       string `json:"realname"`
	Avatar         string `json:"avatar"`
	Nickname       string `json:"nickname"`
	Pid            string `json:"pid"`
}

func TokenGetUser(token string) (userInfo map[string]string) {
	key := "token:" + token
	userInfo = cache.HGetAll(key)
	return
}
