/*
 * @Author: qiuling
 * @Date: 2019-06-20 17:10:45
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-24 15:46:20
 */
package model

// public function MemberRole($user_id int, $role_id int, $client_id string)
// {
// 	$key = "user:memberRole:{$user_id}:{$role_id}:{$client_id}";
// 	return $this->cache->setAutoPrefix(false)->get($key);
// }
import (
	. "artifact/pkg"
	redis "artifact/pkg/cache"
)

var cache = redis.NewCache().SetPrefix("user")

// MemberRoute 获取需要鉴权的路由
func MemberRoute() []string {
	key := "member:route"
	return cache.SMembers(key)
}

func GetRoleIds(route string, method string) []string {
	key := "member:route:" + method + "@" + route
	return cache.SMembers(key)
}

func MemberRole(user_id int64, role_id string, client_id string) (count int64) {
	key := "memberRole:" + Int642String(user_id) + ":" + role_id + ":" + client_id

	// var countStr string
	err := cache.Get(key, &count)
	if err != nil {
		count = 0
	}
	return
}
