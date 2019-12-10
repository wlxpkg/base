package cache

import (
	"encoding/json"
	"fmt"
	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"
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
		log.Err(err)
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
		log.Err(err)
		return
	}
	key = c.prefixKey(key)
	err = clients[c.db].Set(key, string(val), time.Duration(ttl)*time.Second).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// Get cache data
func (c *Cache) Get(key string, structs interface{}) (err error) {
	key = c.prefixKey(key)

	value, err := clients[c.db].Get(key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Err(err)
		}
		return
	}

	str := []byte(value)
	err = json.Unmarshal(str, &structs)
	if err != nil {
		fmt.Printf("err2: \n%#v\n", err)
		log.Err(err)
		return
	}
	return
}

// Del cache data
func (c *Cache) Del(key string) {
	key = c.prefixKey(key)
	err := clients[c.db].Del(key).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// Exists 检测 key 是否存在
func (c *Cache) Exists(key string) (isExists bool) {
	key = c.prefixKey(key)
	value, err := clients[c.db].Exists(key).Result()
	// R(value, "Exists value")

	if err != nil {
		log.Err(err)
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
		log.Err(err)
		return
	}
}

// ExpireAt 设置 key 的过期时间 tm
func (c *Cache) ExpireAt(key string, tm time.Time) {
	key = c.prefixKey(key)
	err := clients[c.db].ExpireAt(key, tm).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// TTL 以秒为单位返回 key 的剩余过期时间
func (c *Cache) TTL(key string) (ttl int64, err error) {
	key = c.prefixKey(key)
	var tm time.Duration
	tm, err = clients[c.db].TTL(key).Result()
	if err != nil {
		log.Err(err)
		return
	}
	ttl = int64(tm / time.Second)

	return
}

/**************************  Incr  Decr ********************************/

// Incr 执行 INCR cmd
func (c *Cache) Incr(key string) {
	key = c.prefixKey(key)
	err := clients[c.db].Incr(key).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// IncrBy 执行 INCRBY cmd
func (c *Cache) IncrBy(key string, incr int64) {
	key = c.prefixKey(key)
	err := clients[c.db].IncrBy(key, incr).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// Decr 执行 DECR cmd
func (c *Cache) Decr(key string) {
	key = c.prefixKey(key)
	err := clients[c.db].Decr(key).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// DecrBy 执行 DECRBY cmd
func (c *Cache) DecrBy(key string, decr int64) {
	key = c.prefixKey(key)
	err := clients[c.db].DecrBy(key, decr).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

/**************************  Hash  ********************************/

// HExists hash 是否存在 field 字段
func (c *Cache) HExists(key, field string) (isExists bool) {
	key = c.prefixKey(key)
	isExists, err := clients[c.db].HExists(key, field).Result()

	if err != nil {
		log.Err(err)
		isExists = false
		return
	}
	return
}

// HIncrBy hash incrBy
func (c *Cache) HIncrBy(key, field string, incr int64) {
	key = c.prefixKey(key)
	err := clients[c.db].HIncrBy(key, field, incr).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// HSet hash
func (c *Cache) HSet(key string, field string, value interface{}) {
	key = c.prefixKey(key)
	err := clients[c.db].HSet(key, field, value).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// HMSet
func (c *Cache) HMSet(key string, data interface{}) {
	value := c.MapData(data)

	key = c.prefixKey(key)
	err := clients[c.db].HMSet(key, value).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

func (c *Cache) HGet(key string, field string) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].HGet(key, field).Result()
	// R(value, "Exists value")

	if err != nil {
		log.Err(err)
		return
	}
	return
}

func (c *Cache) HDel(key string, field string) {
	key = c.prefixKey(key)
	err := clients[c.db].HDel(key, field).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

func (c *Cache) HGetAll(key string) (value map[string]string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].HGetAll(key).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

/******************************* set *********************************/

// SAdd 向集合添加一个或多个成员
func (c *Cache) SAdd(key string, members ...interface{}) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SAdd(key, members...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SCard 获取集合的成员数
func (c *Cache) SCard(key string) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SCard(key).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SDiff 返回给定所有集合的差集
func (c *Cache) SDiff(keys ...string) (value []string) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}

	value, err := clients[c.db].SDiff(keys...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SDiffStore 返回给定所有集合的差集并存储在 destination 中
func (c *Cache) SDiffStore(destination string, keys ...string) (value int64) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}

	value, err := clients[c.db].SDiffStore(destination, keys...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SInter 返回给定所有集合的交集
func (c *Cache) SInter(keys ...string) (value []string) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}

	value, err := clients[c.db].SInter(keys...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SInterStore 返回给定所有集合的交集并存储在 destination 中
func (c *Cache) SInterStore(destination string, keys ...string) (value int64) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}

	value, err := clients[c.db].SInterStore(destination, keys...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SIsMember 判断 member 元素是否是集合 key 的成员
func (c *Cache) SIsMember(key string, member interface{}) (isMember bool) {
	key = c.prefixKey(key)
	isMember, err := clients[c.db].SIsMember(key, member).Result()

	if err != nil {
		log.Err(err)
		isMember = false
		return
	}
	return
}

// SMembers 返回集合中的所有成员
func (c *Cache) SMembers(key string) (value []string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SMembers(key).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SMembersMap 返回集合中的所有成员
func (c *Cache) SMembersMap(key string) (value map[string]struct{}) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SMembersMap(key).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SMove 将 member 元素从 source 集合移动到 destination 集合
func (c *Cache) SMove(source, destination string, member interface{}) (isMove bool) {
	source = c.prefixKey(source)
	destination = c.prefixKey(destination)
	isMove, err := clients[c.db].SMove(source, destination, member).Result()

	if err != nil {
		log.Err(err)
		isMove = false
		return
	}
	return
}

// SPop 移除并返回集合中的一个随机元素
func (c *Cache) SPop(key string) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SPop(key).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SPopN 移除并返回集合中的 n 个随机元素
func (c *Cache) SPopN(key string, count int64) (value []string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SPopN(key, count).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SRandMember 返回集合中的一个随机元素
func (c *Cache) SRandMember(key string) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SRandMember(key).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SRandMemberN 返回集合中的 n 个随机元素
func (c *Cache) SRandMemberN(key string, count int64) (value []string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SRandMemberN(key, count).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SRem 移除集合中一个或多个成员
func (c *Cache) SRem(key string, members ...interface{}) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].SRem(key, members...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SUnion 返回所有给定集合的并集
func (c *Cache) SUnion(keys ...string) (value []string) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}

	value, err := clients[c.db].SUnion(keys...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// SUnionStore 所有给定集合的并集存储在 destination 集合中
func (c *Cache) SUnionStore(destination string, keys ...string) (value int64) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}

	value, err := clients[c.db].SUnionStore(destination, keys...).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}
