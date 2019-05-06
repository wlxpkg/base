/*
 * @Author: qiuling
 * @Date: 2019-05-06 19:00:55
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-06 19:54:32
 */

package test

import (
	. "artifact/pkg"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testKey string = "test:getkey"

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

func TestSetGet(t *testing.T) {
	Set(testKey, student, 100)

	st := Get(testKey, &Student{})
	R(st, "testGet")

	assert.Equal(t, st, student, "student")
}
