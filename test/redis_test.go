/*
 * @Author: qiuling
 * @Date: 2019-05-06 19:00:55
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-24 14:56:56
 */

package test

import (
	. "artifact/pkg"
	redis "artifact/pkg/cache"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Guake bool   `json:"guake"`
	// Classes []string
	Price float32 `json:"price"`
}

var student = Student{
	"Xiao Ming",
	16,
	true,
	// []string{"Math", "English", "Chinese"},
	9.99,
}

var valString string = "test string value"

var cache = redis.NewCache()

func TestSetGet(t *testing.T) {
	testKey := "test:getkey"
	testKeyString := "test:keyString"

	exists := cache.Exists(testKey)
	assert.Equal(t, exists, false, "exists false")

	cache.Set(testKey, student, 2)
	cache.Set(testKeyString, valString, 100)

	exists2 := cache.Exists(testKey)
	assert.Equal(t, exists2, true, "exists true")

	var st Student
	err := cache.Get(testKey, &st)
	R(err, "err")
	R(st, "st")

	var str string
	err2 := cache.Get(testKeyString, &str)
	R(err2, "err2")
	R(str, "testGetStr")

	assert.Equal(t, st, student, "student")
	assert.Equal(t, str, valString, "valString")
}

func TestHSetGet(t *testing.T) {
	testKey := "test:hkey"
	filde := "man"
	hValue := true

	cache.HSet(testKey, filde, hValue)
	cache.Expire(testKey, 100)

	value := cache.HGet(testKey, filde)
	R(value, "hget value")
	assert.Equal(t, hValue, value, "hset hget")

	cache.HMSet(testKey, student)
	cache.HDel(testKey, filde)

	allValue := cache.HGetAll(testKey)
	R(allValue, "hget allValue")

	// var outType Student
	// st := cache.StructData(allValue, &outType)
	// R(st, "st Value")
	// assert.Equal(t, allValue, st, "st value")
}

func TestIdgen(t *testing.T) {
	id, _ := redis.Getid()
	R(id, "Getid")

	ids, _ := redis.Batchid(10)
	R(ids, "Batchid")
}

func TestIncr(t *testing.T) {
	key := "test:incr"
	cache.Del(key)
	cache.Incr(key)
	var incr int
	err := cache.Get(key, &incr)
	cache.Del(key)
	assert.Equal(t, 1, incr, "incr cmd")
	assert.Empty(t, err, "err nil")
}

func TestDecr(t *testing.T) {
	key := "test:decr"
	cache.Del(key)
	cache.Incr(key)
	var decr int
	err := cache.Get(key, &decr)
	cache.Del(key)
	assert.Equal(t, 1, decr, "decr cmd")
	assert.Empty(t, err, "err nil")
}

func TestIncrBy(t *testing.T) {
	key := "test:incrBy"
	cache.Del(key)
	cache.IncrBy(key, 500)
	var incr int
	err := cache.Get(key, &incr)
	cache.Del(key)
	assert.Equal(t, 500, incr, "incrBy cmd")
	assert.Empty(t, err, "err nil")
}

func TestDectBy(t *testing.T) {
	key := "test:decrBy"
	cache.Del(key)
	cache.DecrBy(key, 500)
	var decr int
	err := cache.Get(key, &decr)
	cache.Del(key)
	t.Run("run decrBy success", func(t *testing.T) {
		assert.Equal(t, -500, decr, "decrBy cmd")
		assert.Empty(t, err, "err nil")
	})

	t.Run("run decrBy failed", func(t *testing.T) {
		assert.Equal(t, -100, decr, "decrBy cmd")
		assert.Empty(t, err, "err nil")
	})
}

func TestHExists(t *testing.T) {
	key := "test:hexists"
	cache.Del(key)
	field := "foo"
	hValue := "bar"

	cache.HSet(key, field, hValue)
	cache.Expire(key, 100)
	exists := cache.HExists(key, field)
	t.Run("run hexists success", func(t *testing.T) {
		assert.Equal(t, true, exists, "HExists cmd")
	})
}

func TestHIncr(t *testing.T) {
	key := "test:hincr"
	cache.Del(key)
	field := "foo"

	cache.HIncrBy(key, field, 2)

	result := cache.HGet(key, field)
	t.Run("run hexists success", func(t *testing.T) {
		assert.Equal(t, "2", result, "HIncrBy cmd")
	})
}

func TestSAdd(t *testing.T) {
	key := "test:sadd"
	cache.Del(key)
	Cache.SAdd(key, 1, "2")
	cache.SMembers(key)
	result := cache.SMembers(key)
	t.Run("run sadd success", func(t *testing.T) {
		assert.Equal(t, []string{"1", "2"}, result, "sadd cmd")
	})
}

func TestSCard(t *testing.T) {
	key := "test:scard"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SCard(key)
	t.Run("run scard success", func(t *testing.T) {
		assert.Equal(t, int64(2), result, "scard cmd")
	})
}

func TestSDiff(t *testing.T) {
	key1 := "test:diff:1"
	Cache.SAdd(key1, 1, 2, 3)
	key2 := "test:diff:2"
	Cache.SAdd(key2, 3, 4, 5)
	result := cache.SDiff(key1, key2)
	t.Run("run sdiff success", func(t *testing.T) {
		assert.Equal(t, []string{"1", "2"}, result, "sdiff cmd")
	})
}

func TestSInter(t *testing.T) {
	key1 := "test:inter:1"
	key2 := "test:inter:2"
	cache.Del(key1)
	cache.Del(key2)
	Cache.SAdd(key1, 1, 2, 3)
	Cache.SAdd(key2, 3, 4, 5)
	result := cache.SInter(key1, key2)
	t.Run("run SInter success", func(t *testing.T) {
		assert.Equal(t, []string{"3"}, result, "SInter cmd")
	})
}

func TestSIsMember(t *testing.T) {
	key := "test:SIsMember"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SIsMember(key, 4)
	t.Run("run SIsMember success", func(t *testing.T) {
		assert.Equal(t, true, result, "SIsMember cmd")
	})
}

func TestSPop(t *testing.T) {
	key := "test:SPop"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SPop(key)
	t.Run("run SPop success", func(t *testing.T) {
		assert.Equal(t, "2", result, "SPop cmd")
	})
}

func TestSPopN(t *testing.T) {
	key := "test:SPopN"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SPopN(key, 2)
	t.Run("run SPopN success", func(t *testing.T) {
		assert.Equal(t, []string{"2", "3"}, result, "SPopN cmd")
	})
}

func TestSRandMember(t *testing.T) {
	key := "test:SRandMember"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SRandMember(key)
	t.Run("run SRandMember success", func(t *testing.T) {
		assert.Equal(t, "2", result, "SRandMember cmd")
	})
}

func TestSRandMemberN(t *testing.T) {
	key := "test:SPopN"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SRandMemberN(key, 2)
	t.Run("run SRandMemberN success", func(t *testing.T) {
		assert.Equal(t, []string{"2", "3"}, result, "SRandMemberN cmd")
	})
}

func TestSRem(t *testing.T) {
	key := "test:SRem"
	cache.Del(key)
	Cache.SAdd(key, "2", 3, 4)
	result := cache.SRem(key, 2, 3)

	t.Run("run SRem success", func(t *testing.T) {
		assert.Equal(t, int64(2), result, "SRem cmd")
	})
}

func TestSUnion(t *testing.T) {
	key1 := "test:sunion:1"
	key2 := "test:sunion:2"
	cache.Del(key1)
	cache.Del(key2)
	Cache.SAdd(key1, 1, 2, 3)
	Cache.SAdd(key2, 3, 4, 5)
	result := cache.SUnion(key1, key2)
	t.Run("run SUnion success", func(t *testing.T) {
		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, result, "SUnion cmd")
	})
}

func TestTTL(t *testing.T) {
	key := "test:ttl"
	Cache.Del(key)
	Cache.Set(key, "1", 60)
	tm, _ := Cache.TTL(key)
	t.Run("run ttl success", func(t *testing.T) {
		assert.Equal(t, int64(60), tm, "ttl cmd")
	})
}

func TestExpireAt(t *testing.T) {
	key := "test:expireAt"
	Cache.Del(key)
	Cache.Set(key, "1", 0)
	timeStr := time.Now().Format("2006-01-02")
	t11, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	tomorrow := t11.AddDate(0, 0, 1)

	Cache.ExpireAt(key, tomorrow)
	tm, _ := Cache.TTL(key)
	t.Run("run ttl success", func(t *testing.T) {
		assert.Equal(t, int64(60), tm, "ttl cmd")
	})
}
