/*
 * @Author: qiuling
 * @Date: 2019-07-01 17:12:05
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:38:04
 */
package test

import (
	"testing"
	"time"

	. "github.com/wlxpkg/base"
	"github.com/wlxpkg/base/amqp"
	"github.com/wlxpkg/base/log"
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
	data := make(map[string]string)

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
