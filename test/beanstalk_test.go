/*
 * @Author: qiuling
 * @Date: 2019-06-28 19:13:57
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-01 10:57:29
 */
package test

import (
	. "artifact/pkg"
	"artifact/pkg/beanstalk"
	"artifact/pkg/log"
	"testing"
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

	jobID, err := publish()

	R(jobID, "job")
	R(err, "err")

}

func publish() (uint64, error) {

	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["type"] = "99"
	data["is_default"] = "0"

	jobID, err := pPool.Publish("/test/publist", data, 15)

	R(jobID, "job")
	R(err, "err")
	return jobID, err
}

func TestConsumer(t *testing.T) {
	beanstalk.NewConsumer("test", func(jsonStr string) (bool, error) {
		R(jsonStr, "jsonStr")
		json, _ := JsonDecode(jsonStr)

		topic := json["topic"].(string)
		R(topic, "topic")

		message := json["message"].(map[string]interface{})
		R(message, "message")

		msgName := message["name"].(string)
		R(msgName, "msgName")

		return true, nil
	})
}
