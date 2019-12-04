/*
 * @Author: qiuling
 * @Date: 2019-07-01 20:13:06
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-02 17:19:25
 */
package amqp

import (
	. "git.wlx/zwyd/pkg"
	"git.wlx/zwyd/pkg/log"

	"github.com/streadway/amqp"
)

// Receiver 观察者模式需要的接口
// 观察者用于接收指定的queue到来的数据
type Receiver interface {
	QueueName() string            // 获取接收者需要监听的队列
	AutoDelete() bool             // 队列不存在是否自动删除
	OnError(error)                // 处理遇到的错误，当RabbitMQ对象发生了错误，他需要告诉接收者处理错误
	OnReceive(amqp.Delivery) bool // 处理收到的消息, 这里需要告知RabbitMQ对象消息是否处理成功
}

type Receive struct {
	queue      string
	autoDelete bool
}

func NewReceiver(queue string, autoDelete bool) *Receive {
	r := &Receive{
		queue,
		autoDelete,
	}

	return r
}

func (r *Receive) QueueName() string {
	return r.queue
}

func (r *Receive) AutoDelete() bool {
	return r.autoDelete
}

func (r *Receive) OnError(err error) {
	log.Warn(err)
}

// OnReceive 消息处理, 此处为demo
func (r *Receive) OnReceive(msg amqp.Delivery) (ok bool) {
	R(Byte2String(msg.Body), "收到 AMQP消息")

	_ = msg.Ack(false)
	return true
}
