/*
 * @Author: qiuling
 * @Date: 2019-05-06 19:00:55
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-17 17:55:13
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

var student = &Student{
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

	st := cache.Get(testKey, &Student{})
	R(st, "testGet")

	str := cache.Get(testKeyString, "")
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
