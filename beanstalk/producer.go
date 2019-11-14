/*
 * @Author: qiuling
 * @Date: 2019-06-28 15:38:14
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-11-11 10:31:08
 */
package beanstalk

import (
	. "artifact/pkg"
	. "artifact/pkg/config"
	"artifact/pkg/log"
	"strings"
	"time"

	"github.com/beanstalkd/go-beanstalk"
	"github.com/silenceper/pool"
)

type Producer struct {
	pool pool.Pool
}

func NewProducer() (producer *Producer, err error) {
	producer = &Producer{}

	//factory 创建连接的方法
	factory := func() (interface{}, error) {
		return beanstalk.Dial("tcp", Config.Beanstalk.Host+":"+Config.Beanstalk.Port)
	}
	//close 关闭连接的方法
	close := func(v interface{}) error {
		return v.(*beanstalk.Conn).Close()
	}
	//ping 检测连接的方法
	//ping := func(v interface{}) error { return nil }

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &pool.PoolConfig{
		InitialCap: 5,
		MaxCap:     30,
		Factory:    factory,
		Close:      close,
		//Ping:       ping,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}
	pool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		R(err, "err")
	}

	producer.pool = pool

	return
}

// Publish 发送延迟消息
func (producer *Producer) Publish(
	topic string,
	message interface{},
	delay int64,
) (uint64, error) {

	messageStr, _ := JsonEncode(message)
	R("发送延迟消息: "+topic+", 消息内容: "+Byte2String(messageStr)+", 延迟:"+Int642String(delay)+" 秒", "")

	routing := strings.Split(topic, "/")
	tube := routing[1]

	msg := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}

	messageData, err := JsonEncode(msg)
	if err != nil {
		log.Info(err)
		return 0, err
	}

	puter, err := producer.pool.Get()

	if err != nil {
		log.Info(err)
		return 0, err
	}

	conn := puter.(*beanstalk.Conn)

	// 修改 tube
	if conn.Tube.Name != tube {
		conn.Tube = beanstalk.Tube{Conn: conn, Name: tube}
	}

	id, err := conn.Put(messageData, 1, time.Duration(delay)*time.Second, 30*time.Second)

	//将连接放回连接池中
	_ = producer.pool.Put(puter)

	return id, err
}
