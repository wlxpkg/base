/*
 * @Author: qiuling
 * @Date: 2019-06-17 15:57:14
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:10:17
 */
package cache

import (
	"errors"
	"strings"

	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"

	"github.com/go-redis/redis"
)

// type Idgenerator struct {
// }

var idgeneratorClient *redis.Client

func init() {
	idgeneratorClient = redis.NewClient(&redis.Options{
		Addr:     Config.Idgenerator.Host + ":" + Config.Idgenerator.Port,
		Password: "",
		DB:       0,
		PoolSize: 100,
	})
}

/* func NewCache() (cache *Cache) {
	cache = &Cache{
		autoPrefix: true,
		prefix:     Config.Redis.Prefix,
	}
	return
} */

func Getid() (value string, err error) {
	value, err = idgeneratorClient.Do("getid").String()

	if err != nil {
		log.Warn(err)
		return
	}

	if strings.Contains(value, "ERROR") {
		log.Warn(err)
		err = errors.New("ERR_IDGEN_FAIL")
		return
	}

	return
}

func Batchid(num int) (value string, err error) {
	value, err = idgeneratorClient.Do("batchid", num).String()

	if err != nil {
		log.Warn(err)
		return
	}

	if strings.Contains(value, "ERROR") {
		log.Warn(err)
		err = errors.New("ERR_IDGEN_FAIL")
		return
	}

	return
}
