/*
 * @Author: qiuling
 * @Date: 2019-06-28 19:13:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-28 23:43:44
 */
package test

import (
	. "artifact/pkg"
	"artifact/pkg/beanstalk"
	"artifact/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var data = make(map[string]string)
func TestPublish(t *testing.T) {
	pPool := beanstalk.NewProducer()

	// R(pPool, "pooooo")
	// err := pPool.Conn()
	// R(err, "err")
	// R(pPool, "p1")

	pool, err := pPool.NewPool()
	if err != nil {
		log.Err("Unable to create beanstalk producer pool: " + err.Error())
	}
	defer pool.Stop()

	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["permission_id"] = "4,10,2,1,16,19,23,27,31,35,39,40,36,32,28,24,20,17,18,21,22,26,25,29,30,33,34,37,38,41,42"
	data["type"] = "99"
	data["is_default"] = "0"

	tube, message, param, _ := pPool.MakeData("/test/publist", data, 15)

	jobID, err := pool.Put(tube, message, param)

	R(jobID, "job")
	R(err, "err")

	assert.Equal(t, 1, 1, "TestLog")
}
