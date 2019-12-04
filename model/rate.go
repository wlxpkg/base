/*
 * @Author: qiuling
 * @Date: 2019-09-19 10:27:59
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-09-19 17:57:26
 */
package model

import (
	. "git.wlx/zwyd/pkg"
	redis "git.wlx/zwyd/pkg/cache"
	. "git.wlx/zwyd/pkg/config"
)

const cacheTime = 3

var whitelist = Config.Rate.Whitelist
var blacklist = Config.Rate.Blacklist
var longTime = Config.Rate.LongTime
var shortTime = Config.Rate.ShortTime

var rateCache = redis.NewCache().SetDB("rate").SetPrefix("rate")

func RateCheck(jwt string, method string, path string) bool {
	// 排除 黑名单之外的 GET 请求
	uri := method + "@" + path
	if method == "GET" && !InArray(uri, blacklist) {
		return true
	}
	// 白名单
	if InArray(uri, whitelist) {
		return true
	}
	time := cacheTime
	// 低频
	if InArray(uri, longTime) {
		time = 30
	}
	// 高频
	if InArray(uri, shortTime) {
		time = 1
	}

	jwt = Md5(jwt)
	key := jwt + ":" + uri

	exists := rateCache.Exists(key)

	if !exists {
		rateCache.Set(key, 1, time)
		return true
	}

	return false
}
