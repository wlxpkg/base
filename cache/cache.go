package cache

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"

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
func (c *Cache) Incr(key string) int64 {
	key = c.prefixKey(key)
	return clients[c.db].Incr(key).Val()
}

// IncrBy 执行 INCRBY cmd
func (c *Cache) IncrBy(key string, incr int64) int64 {
	key = c.prefixKey(key)
	return clients[c.db].IncrBy(key, incr).Val()
}

// Decr 执行 DECR cmd
func (c *Cache) Decr(key string) int64 {
	key = c.prefixKey(key)
	return clients[c.db].Decr(key).Val()
}

// DecrBy 执行 DECRBY cmd
func (c *Cache) DecrBy(key string, decr int64) int64 {
	key = c.prefixKey(key)
	return clients[c.db].DecrBy(key, decr).Val()
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
func (c *Cache) HIncrBy(key, field string, incr int64) int64 {
	key = c.prefixKey(key)
	return clients[c.db].HIncrBy(key, field, incr).Val()
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

// HMSet 同时将多个 field-value (字段-值)对设置到哈希表中
func (c *Cache) HMSet(key string, data interface{}) {
	value := c.MapData(data)

	key = c.prefixKey(key)
	err := clients[c.db].HMSet(key, value).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// HGet 返回哈希表中指定字段的值
func (c *Cache) HGet(key string, field string) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].HGet(key, field).Result()

	if err != nil {
		log.Err(err)
		return
	}
	return
}

// HDel 删除哈希表 key 中的一个或多个指定字段
func (c *Cache) HDel(key string, field string) {
	key = c.prefixKey(key)
	err := clients[c.db].HDel(key, field).Err()
	if err != nil {
		log.Err(err)
		return
	}
}

// HGetAll 返回哈希表中，所有的字段和
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

/******************************* list *********************************/

// BLPop 移出并获取列表的第一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
func (c *Cache) BLPop(timeout time.Duration, keys ...string) (value []string) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}
	value, err := clients[c.db].BLPop(timeout, keys...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// BRPop 移出并获取列表的最后一个元素， 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止
func (c *Cache) BRPop(timeout time.Duration, keys ...string) (value []string) {
	for i, key := range keys {
		keys[i] = c.prefixKey(key)
	}
	value, err := clients[c.db].BRPop(timeout, keys...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// BRPopLPush 从列表中弹出一个值，将弹出的元素插入到另外一个列表中并返回它;
// 如果列表没有元素会阻塞列表直到等待超时或发现可弹出元素为止。
func (c *Cache) BRPopLPush(source, destination string, timeout time.Duration) (value string) {
	source = c.prefixKey(source)
	destination = c.prefixKey(destination)
	value, err := clients[c.db].BRPopLPush(source, destination, timeout).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LIndex 通过索引获取列表中的元素
func (c *Cache) LIndex(key string, index int64) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LIndex(key, index).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LInsert 在列表的元素前或者后插入元素
func (c *Cache) LInsert(key, op string, pivot, value interface{}) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LInsert(key, op, pivot, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LInsertBefore 在列表的元素前插入元素
func (c *Cache) LInsertBefore(key string, pivot, value interface{}) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].LInsertBefore(key, pivot, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LInsertAfter 在列表的元素后插入元素
func (c *Cache) LInsertAfter(key string, pivot, value interface{}) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].LInsertAfter(key, pivot, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LLen 获取列表长度
func (c *Cache) LLen(key string) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LLen(key).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LPop 移出并获取列表的第一个元素
func (c *Cache) LPop(key string) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LLen(key).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LPush 将一个或多个值插入到列表头部
func (c *Cache) LPush(key string, values ...interface{}) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LPush(key, values...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LPushX 将一个值插入到已存在的列表头部
func (c *Cache) LPushX(key string, value interface{}) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].LPushX(key, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LRange 获取列表指定范围内的元素 区间以偏移量 START 和 STOP 指定
func (c *Cache) LRange(key string, start, stop int64) (value []string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LRange(key, start, stop).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LRem 根据参数 COUNT 的值，移除列表中与参数 VALUE 相等的元素
func (c *Cache) LRem(key string, count int64, value interface{}) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].LRem(key, count, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// LSet 通过索引来设置元素的值
func (c *Cache) LSet(key string, index int64, value interface{}) (res string) {
	key = c.prefixKey(key)
	res, err := clients[c.db].LSet(key, index, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return

}

// LTrim 让列表只保留指定区间内的元素，不在指定区间之内的元素都将被删除
func (c *Cache) LTrim(key string, start, stop int64) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].LTrim(key, start, stop).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// RPop 移除列表的最后一个元素，返回值为移除的元素
func (c *Cache) RPop(key string) (value string) {
	key = c.prefixKey(key)
	value, err := clients[c.db].RPop(key).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// RPopLPush 移除列表的最后一个元素，并将该元素添加到另一个列表并返回
func (c *Cache) RPopLPush(source, destination string) (value string) {
	source = c.prefixKey(source)
	destination = c.prefixKey(destination)
	value, err := clients[c.db].RPopLPush(source, destination).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// RPush 将一个或多个值插入到列表的尾部
func (c *Cache) RPush(key string, values ...interface{}) (value int64) {
	key = c.prefixKey(key)
	value, err := clients[c.db].RPush(key, values...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// RPushX 将一个值插入到列表的尾部
func (c *Cache) RPushX(key string, value interface{}) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].RPushX(key, value).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

/******************************* sorted set *********************************/

// ZAdd 将一个或多个成员元素及其分数值加入到有序集当中
func (c *Cache) ZAdd(key string, members ...redis.Z) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZAdd(key, members...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZAddNX 将一个或多个成员元素及其分数值加入到有序集当中
func (c *Cache) ZAddNX(key string, members ...redis.Z) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZAddNX(key, members...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZIncr 对有序集合中指定成员的分数加上1
func (c *Cache) ZIncr(key string, member redis.Z) (res float64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZIncr(key, member).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZIncrNX 对有序集合中指定成员的分数加上1
func (c *Cache) ZIncrNX(key string, member redis.Z) (res float64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZIncrNX(key, member).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZIncrBy 对有序集合中指定成员的分数加上increment
func (c *Cache) ZIncrBy(key string, increment float64, member string) (res float64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZIncrBy(key, increment, member).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZCard 计算集合中元素的数量
func (c *Cache) ZCard(key string) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZCard(key).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZScore 返回有序集中，成员的分数值
func (c *Cache) ZScore(key, member string) (res float64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZScore(key, member).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZRange 返回有序集中，按分数值递增返回指定区间内的成员
func (c *Cache) ZRange(key string, start, stop int64) (res []string) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZRange(key, start, stop).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZRevRange 返回有序集中，按分数值递减返回指定区间内的成员
func (c *Cache) ZRevRange(key string, start, stop int64) (res []string) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZRevRange(key, start, stop).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}

// ZRem 移除有序集合中的一个或多个成员
func (c *Cache) ZRem(key string, members ...interface{}) (res int64) {
	key = c.prefixKey(key)
	res, err := clients[c.db].ZRem(key, members...).Result()
	if err != nil {
		log.Err(err)
		return
	}
	return
}
