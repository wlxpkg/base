/*
 * @Author: qiuling
 * @Date: 2019-06-20 17:10:45
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-10 18:40:13
 */
package model

// public function MemberRole($user_id int, $role_id int, $client_id string)
// {
// 	$key = "user:memberRole:{$user_id}:{$role_id}:{$client_id}";
// 	return $this->cache->setAutoPrefix(false)->get($key);
// }
import (
	. "github.com/wlxpkg/base"
	redis "github.com/wlxpkg/base/cache"
)

var cache = redis.NewCache().SetPrefix("user")

// GuestRoute 游客路由白名单
func GuestRoute() []string {
	key := "guest:route"
	return cache.SMembers(key)
}

// MemberRoute 获取需要鉴权的路由
func MemberRoute() []string {
	key := "member:route"
	return cache.SMembers(key)
}

/* func GetRoleIds(route string, method string) []string {
	key := "member:route:" + method + "@" + route
	return cache.SMembers(key)
} */

func MemberTime(user_id int64, client_id string) (expireAt string) {
	key := "memberTime:" + Int642String(user_id) + ":" + client_id

	// var countStr string
	err := cache.Get(key, &expireAt)
	if err != nil {
		expireAt = ""
	}
	return
}
