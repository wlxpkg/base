/*
 * @Author: qiuling
 * @Date: 2019-07-01 17:09:59
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-08-07 10:17:59
 */
package amqp

import (
	"fmt"
	. "git.wlx/zwyd/pkg"
	"git.wlx/zwyd/pkg/log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// Consumer 用于管理和维护rabbitmq的对象
type Consumer struct {
	wg sync.WaitGroup

	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string // exchange的名称
	exchangeType string // exchange的类型
	receivers    []Receiver
}

func NewConsumer() *Consumer {
	return &Consumer{
		exchangeName: exchange,
		exchangeType: "direct",
	}
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
func (c *Consumer) RegisterReceiver(receiver Receiver) {
	c.receivers = append(c.receivers, receiver)
}

// Start 启动Rabbitmq的客户端
func (c *Consumer) Start() {
	for {
		c.run()

		R(time.Now().Format("2006-01-02 15:04:05"), " amqp 断开连接")
		// 一旦连接断开，那么需要隔一段时间去重连
		// 这里最好有一个时间间隔
		time.Sleep(3 * time.Second)
	}
}

// run 开始获取连接并初始化相关操作
func (c *Consumer) run() {
	/* if !config.Global.RabbitMQ.Refresh() {
		log.Err("rabbit刷新连接失败，将要重连: %s", config.Global.RabbitMQ.URL)
		return
	} */

	// 获取新的channel对象
	c.conn, c.channel, _ = Conn()
	R("连接成功", "AMQP 消费队列 ")

	// 初始化Exchange
	_ = PrepareExchange(c.channel)

	for _, receiver := range c.receivers {
		c.wg.Add(1)
		go c.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
	}

	c.wg.Wait()

	log.Err("所有处理queue的任务都意外退出了")

	// 理论上c.run()在程序的执行过程中是不会结束的
	// 一旦结束就说明所有的接收者都退出了，那么意味着程序与rabbitmq的连接断开
	// 那么则需要重新连接，这里尝试销毁当前连接
	// config.Global.RabbitMQ.Distory()
	c.distory()
}

// Listen 监听指定路由发来的消息
// 这里需要针对每一个接收者启动一个goroutine来执行listen
// 该方法负责从每一个接收者监听的队列中获取数据，并负责重试
func (c *Consumer) listen(receiver Receiver) {
	defer c.wg.Done()

	// 这里获取每个接收者需要监听的队列和路由
	queueName := receiver.QueueName()

	// 申明Queue
	_, err := c.channel.QueueDeclare(
		queueName,             // name
		false,                 // durable
		receiver.AutoDelete(), // delete when usused
		false,                 // exclusive(排他性队列)
		false,                 // no-wait
		nil,                   // arguments
	)
	// R(queue, "queue")
	if nil != err {
		receiver.OnError(fmt.Errorf("初始化队列 %s 失败: %s", queueName, err.Error()))
	}

	// 将Queue绑定到Exchange上去
	err = c.channel.QueueBind(
		queueName,      // queue name
		queueName,      // routing key
		c.exchangeName, // exchange
		false,          // no-wait
		nil,
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("绑定队列 [%s - %s] 到交换机失败: %s", queueName, queueName, err.Error()))
	}

	// R(queueName, "绑定队列")

	// 获取消费通道
	_ = c.channel.Qos(2, 0, false)
	msgs, err := c.channel.Consume(
		queueName,               // queue
		"consumer-"+RandStr(10), // consumer
		false,                   // auto-ack
		false,                   // exclusive
		false,                   // no-local
		false,                   // no-wait
		nil,                     // args
	)
	if nil != err {
		receiver.OnError(fmt.Errorf("获取队列 %s 的消费通道失败: %s", queueName, err.Error()))
	}

	// 使用callback消费数据
	for msg := range msgs {
		// R(msg, "AMQP 收到消息")

		// 当接收者消息处理失败的时候，
		// 比如网络问题导致的数据库连接失败，redis连接失败等等这种
		// 通过重试可以成功的操作，那么这个时候是需要重试的
		// 直到数据处理成功后再返回，然后才会回复rabbitmq ack
		for !receiver.OnReceive(msg) {
			log.Warn("receiver 数据处理失败，将要重试")
			time.Sleep(1 * time.Second)
		}

		// 确认收到本条消息, multiple必须为false
		// _ = msg.Ack(false)
	}
}

// Ack 确认消息
func Ack(msg amqp.Delivery) {
	_ = msg.Ack(false)
}

// Reject 拒绝消息
// requeue 是否重新入队列
func Reject(msg amqp.Delivery, requeue bool) {
	_ = msg.Reject(requeue)
}

func (c *Consumer) distory() {
	go c.conn.Close()
}
