/*
 * @Author: qiuling
 * @Date: 2019-06-28 19:13:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:38:14
 */
package test

import (
	"testing"
	"time"

	. "github.com/wlxpkg/base"
	"github.com/wlxpkg/base/beanstalk"
	"github.com/wlxpkg/base/log"
)

var producer = newProducer()

func newProducer() *beanstalk.Producer {
	pool, err := beanstalk.NewProducer()

	if err != nil {
		log.Err("Unable to create beanstalk producer pool: " + err.Error())
	}

	return pool
}

// var data = make(map[string]string)
func TestPublish(t *testing.T) {

	// defer pPool.Stop()

	consumer := beanstalk.NewConsumer("test")
	receiver := beanstalk.NewReceiver()

	consumer.RegisterReceiver(receiver)

	for i := 1; i <= 10; i++ {
		// 延迟执行
		time.AfterFunc(time.Duration(i*10)*time.Second, func() {
			publish()
		})
	}

	consumer.Start()
}

func publish() {
	data := make(map[string]string)

	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["type"] = "99"
	data["is_default"] = "0"

	jobID, err := producer.Publish("/test/publist", data, 15)

	R(jobID, "job")
	R(err, "err")
}
