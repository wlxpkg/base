/*
 * @Author: qiuling
 * @Date: 2019-07-01 11:11:40
 * @Last Modified by: qiuling
 * @Last Modified time: 2020-02-05 21:20:03
 */
package amqp

import (
	"strings"

	. "github.com/wlxpkg/base"
	"github.com/wlxpkg/base/log"

	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewProducer() *Producer {
	p := &Producer{}
	return p
}

// Publish amqp 发送消息
func (p *Producer) Publish(topic string, message interface{}) error {
	messageStr, _ := JsonEncode(message)

	R("topic: "+topic+" message:"+Byte2String(messageStr), "AMQP 发送消息: ")

	routing := strings.Split(topic, "/")
	queue := routing[1]

	msg := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}

	messageData, _ := JsonEncode(msg)

	amqpmsg := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         messageData,
		DeliveryMode: 2,
	}

	if p.channel == nil || p.conn.IsClosed() {
		p.connect(queue)
	}

	return p.channel.Publish(exchange, queue, false, false, amqpmsg)
}

func (p *Producer) connect(queueName string) {
	p.conn, p.channel, _ = Conn()

	err := PrepareExchange(p.channel)
	if err != nil {
		log.Err("AMQP 交换机绑定失败: " + err.Error())
		return
	}
	R("连接成功!", "AMQP 生产队列")
}
