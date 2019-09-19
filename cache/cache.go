package cache

import (
	. "artifact/pkg/config"
	"artifact/pkg/log"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

type Cache struct {
	autoPrefix bool
	prefix     string
	db         string
}

// var client *redis.Client
var clients = make(map[string]*redis.Client)

func init() {
	clients["cache"] = redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + ":" + Config.Redis.Port,
		Password: Config.Redis.Password,
		DB:       Config.Redis.Select,
		PoolSize: 100,
	})

	clients["rate"] = redis.NewClient(&redis.Options{
		Addr:     Config.RateRedis.Host + ":" + Config.RateRedis.Port,
		Password: Config.RateRedis.Password,
		DB:       Config.RateRedis.Select,
		PoolSize: 100,
	})
}

func NewCache() (cache *Cache) {
	cache = new(Cache)
	cache.autoPrefix = true
	cache.prefix = Config.Redis.Prefix
	cache.db = "cache"
	return
}

func (c *Cache) SetAutoPrefix(auto bool) *Cache {
	c.autoPrefix = auto
	return c
}

func (c *Cache) SetPrefix(prefix string) *Cache {
	c.prefix = prefix
	return c
}

func (c *Cache) SetDB(db string) *Cache {
	c.db = db
	return c
}

func (c *Cache) prefixKey(key string) string {
	if c.autoPrefix && !strings.Contains(key, c.prefix) {
		return c.prefix + ":" + key
	}
	c.SetAutoPrefix(true)
	return key
}

// MapData 将数据转为 map 结构
func (c *Cache) MapData(data interface{}) map[string]interface{} {
	mdata := make(map[string]interface{})
	j, _ := json.Marshal(data)

	err := json.Unmarshal(j, &mdata)
	if err != nil {
		log.Warn(err)
		return mdata
	}
	return mdata
}

/**********************************************************************/

// Set cache data
// ttl Second
func (c *Cache) Set(key string, value interface{}, ttl int) {
	val, err := json.Marshal(value)
	if err != nil {
		log.Warn(err)
		return
	}
	key = c.prefixKey(key)
	err = clients[c.db].Set(key, string(val), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Warn(err)
		return
	}
}

// Get cache data
func (c *Cache) Get(key string, structs interface{}) (err error) {
	key = c.prefixKey(key)

	value, err := clients[c.db].Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Warn(err)
		}
		return
	}

	str := []byte(value)
	err = json.Unmarshal(str, &structs)
	if err != nil {
		fmt.Printf("err2: \n%#v\n", err)
		log.Warn(err)
		return
	}
	return
}

// Del cache data
func (c *Cache) Del(key string) {
	key = c.prefixKey(key)
	err := clients[c.db].Del(key).Err()
	if err != nil {
		log.Warn(err)
		return
	}
}

// Exists 检测 key 是否存在
func (c *Cache) Exists(key string) (isExists bool) {
	key = c.prefixKey(key)
	value, err := clients[c.db].Exists(key).Result()
	// R(value, "Exists value")

	if err != nil {
		log.Warn(err)
		return false
	}

	if value > 0 {
		return true
	}
	return false
}

// Expire 设置 key 的有效期为 ttl 秒
func (c *Cache) Expire(key string, ttl int) {
	key = c.prefixKey(key)
	err := clients[c.db].Expire(key, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Warn(err)
		return
	}
}

/**************************  Hash  ********************************/

// HSet
func (c *Cache) HSet(key string, field string, value interface{}) {
	key = c.prefixKey(key)
	err := clients[c.db].HSet(key, field, value).Err()
	if err != nil {
		log.Warn(err)
		return
	}
}

// HMSet
func (c *Cache) HMSet(key string, data interface{}) {
	value := c.MapData(data)

	key = c.prefixKey(key)
	err := clients[c.db].HMSet(key, value).Err()
	if err != nil {
		log.Warn(err)
		return
	}
}

func (c *Cache) HGet(key string, field string) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].HGet(key, field).Result()
	// R(value, "Exists value")

	if err != nil {
		log.Warn(err)
		return
	}
	return
}

func (c *Cache) HDel(key string, field string) {
	key = c.prefixKey(key)
	err := clients[c.db].HDel(key, field).Err()
	if err != nil {
		log.Warn(err)
		return
	}
}

func (c *Cache) HGetAll(key string) (value map[string]string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].HGetAll(key).Result()

	if err != nil {
		log.Warn(err)
		return
	}
	return
}

/******************************* set *********************************/

func (c *Cache) SMembers(key string) (value []string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SMembers(key).Result()

	if err != nil {
		log.Warn(err)
		return
	}
	return
}
