/*
 * @Author: qiuling
 * @Date: 2019-07-01 17:49:56
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-01 22:06:01
 */
package amqp

import (
	// . "artifact/pkg"
	. "artifact/pkg/config"
	"artifact/pkg/log"
	"strings"

	"github.com/streadway/amqp"
)

const exchange = "artifact"
const exchangeType = "direct"

func Conn() (conn *amqp.Connection, channel *amqp.Channel, err error) {
	vhostArr := strings.Split(Config.Amqp.Vhost, "/")
	vhost := "%2f" + vhostArr[1]
	url := "amqp://" + Config.Amqp.User + ":" + Config.Amqp.Pass + "@" + Config.Amqp.Host + ":" + Config.Amqp.Port + "/" + vhost

	conn, err = amqp.Dial(url)

	if err != nil {
		log.Err("AMQP 连接失败: " + err.Error())
		return
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Err("AMQP 连接失败: " + err.Error())
		return
	}

	// R(channel, "channel")
	return
}

// prepareExchange 准备rabbitmq的Exchange
func PrepareExchange(channel *amqp.Channel) (err error) {
	// 申明Exchange
	return channel.ExchangeDeclare(
		exchange,     // exchange
		exchangeType, // type
		true,         // durable
		false,        // autoDelete
		false,        // internal
		false,        // noWait
		nil,          // args
	)
}
