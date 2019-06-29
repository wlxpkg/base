/*
 * @Author: qiuling
 * @Date: 2019-06-28 19:13:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-29 16:13:49
 */
package test

import (
	. "artifact/pkg"
	"artifact/pkg/beanstalk"
	"artifact/pkg/log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pPool *beanstalk.Producer

// var data = make(map[string]string)
func TestPublish(t *testing.T) {
	var err error
	pPool, err = beanstalk.NewProducerPool()

	if err != nil {
		log.Err("Unable to create beanstalk producer pool: " + err.Error())
	}
	defer pPool.Stop()

	jobID, err := publish()

	R(jobID, "job")
	R(err, "err")

	assert.Equal(t, 1, 1, "TestLog")
}

func publish() (uint64, error) {

	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["permission_id"] = "4,10,2,1,16,19,23,27,31,35,39,40,36,32,28,24,20,17,18,21,22,26,25,29,30,33,34,37,38,41,42"
	data["type"] = "99"
	data["is_default"] = "0"

	jobID, err := pPool.Publish("/test/publist", data, 15)

	R(jobID, "job")
	R(err, "err")
	return jobID, err
}
