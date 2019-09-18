/*
 * @Author: qiuling
 * @Date: 2019-07-01 17:12:05
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-09-17 16:37:37
 */
package test

import (
	. "artifact/pkg"
	"artifact/pkg/amqp"
	"artifact/pkg/log"
	"testing"
	"time"
)

func TestAmqpPublish(t *testing.T) {
	pub := amqp.NewProducer()
	csm := amqp.NewConsumer()

	recv := amqp.NewReceiver("test", true)

	csm.RegisterReceiver(recv)

	for i := 1; i <= 3; i++ {
		// 延迟执行
		time.AfterFunc(time.Duration(i*5)*time.Second, func() {
			amqpPublish(pub)
		})
	}

	csm.Start()
}

func amqpPublish(pub *amqp.Producer) {
	data["name"] = "测试角色"
	data["slug"] = "customer"
	data["type"] = "99"
	data["is_default"] = "0"

	err := pub.Publish("/test/publist", data)

	if err != nil {
		log.Warn(err)
		R("", "发送失败")
	} else {
		R("发送成功", "")
	}

}
