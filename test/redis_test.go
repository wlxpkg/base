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
