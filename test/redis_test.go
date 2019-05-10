/*
 * @Author: qiuling
 * @Date: 2019-05-06 19:00:55
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-10 16:45:35
 */

package test

import (
	. "artifact/pkg"
	"artifact/pkg/cache"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testKey string = "test:getkey"
var testKeyString string = "test:keyString"

type Student struct {
	Name    string
	Age     int
	Guake   bool
	Classes []string
	Price   float32
}

var student = &Student{
	"Xiao Ming",
	16,
	true,
	[]string{"Math", "English", "Chinese"},
	9.99,
}

var valString string = "test string value"

func TestSetGet(t *testing.T) {
	cache.Set(testKey, student, 100)
	cache.Set(testKeyString, valString, 100)

	st := cache.Get(testKey, &Student{})
	R(st, "testGet")

	str := cache.Get(testKeyString, "")
	R(str, "testGetStr")

	assert.Equal(t, st, student, "student")
	assert.Equal(t, str, valString, "valString")
}
