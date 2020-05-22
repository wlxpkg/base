/*
 * @Author: qiuling
 * @Date: 2019-07-01 17:49:56
 * @Last Modified by: qiuling
 * @Last Modified time: 2020-05-22 22:23:34
 */
package amqp

import (
	"strings"

	. "github.com/wlxpkg/base/config"
	"github.com/wlxpkg/base/log"

	"github.com/streadway/amqp"
)

const exchangeType = "direct"

var exchange = Config.Amqp.Exchange

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
