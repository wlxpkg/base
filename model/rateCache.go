/*
 * @Author: qiuling
 * @Date: 2019-09-19 10:27:59
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-09-19 11:48:56
 */
package model

import (
	. "artifact/pkg"
	redis "artifact/pkg/cache"
	. "artifact/pkg/config"
)

const cacheTime = 3

var whitelist = Config.Rate.Whitelist
var longTime = Config.Rate.LongTime
var shortTime = Config.Rate.ShortTime

var rateCache = redis.NewCache().SetDB("rate").SetPrefix("rate")

func RateCheck(jwt string, path string) bool {
	// 白名单
	if InArray(path, whitelist) {
		return true
	}
	time := cacheTime
	// 低频
	if InArray(path, longTime) {
		time = 30
	}
	// 高频
	if InArray(path, shortTime) {
		time = 1
	}

	jwt = Md5(jwt)
	key := jwt + ":" + path

	exists := rateCache.Exists(key)

	if !exists {
		rateCache.Set(key, 1, time)
		return true
	}

	return false
}
