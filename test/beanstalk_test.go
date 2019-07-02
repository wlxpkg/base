/*
 * @Author: qiuling
 * @Date: 2019-06-28 19:13:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-02 17:29:48
 */
package test

import (
	. "artifact/pkg"
	"artifact/pkg/beanstalk"
	"artifact/pkg/log"
	"testing"
	"time"
)

var pPool *beanstalk.ProducerPool

// var data = make(map[string]string)
func TestPublish(t *testing.T) {
	var err error
	pPool, err = beanstalk.NewProducerPool()

	if err != nil {
		log.Err("Unable to create beanstalk producer pool: " + err.Error())
	}
	defer pPool.Stop()

	consumer := beanstalk.NewConsumer("test")
	receiver := beanstalk.NewReceiver()

	consumer.RegisterReceiver(receiver)

	for i := 1; i <= 3; i++ {
		// 延迟执行
		time.AfterFunc(time.Duration(i*5)*time.Second, func() {
			publish()
		})
	}

	consumer.Start()
}

func publish() {

	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["type"] = "99"
	data["is_default"] = "0"

	jobID, err := pPool.Publish("/test/publist", data, 15)

	R(jobID, "job")
	R(err, "err")
}
