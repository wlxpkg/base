/*
 * @Author: qiuling
 * @Date: 2019-06-28 15:38:14
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-29 17:17:37
 */
package beanstalk

import (
	. "artifact/pkg"
	. "artifact/pkg/config"
	graylog "artifact/pkg/log"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	bt "github.com/prep/beanstalk"
)

type ProducerPool struct {
	producers []*bt.Producer
	putC      chan *bt.Put
	putTokens chan *bt.Put
	stopOnce  sync.Once
}

func NewProducerPool() (*ProducerPool, error) {
	pool := &ProducerPool{putC: make(chan *bt.Put)}

	link := "beanstalk://" + Config.Beanstalk.Host + ":" + Config.Beanstalk.Port

	var urls []string
	urls = append(urls, link, link, link)

	options := &bt.Options{
		// ReserveTimeout defines how long a beanstalk reserve command should wait
		// before it should timeout. The default and minimum value is 1 second.
		ReserveTimeout: 3 * time.Second,
		// ReconnectTimeout defines how long a producer or consumer should wait
		// between reconnect attempts. The default is 3 seconds, with a minimum of 1
		// second.
		ReconnectTimeout: 3 * time.Second,
		// ReadWriteTimeout defines how long each read or write operation is  allowed
		// to block until the connection is considered broken. The default is
		// disabled and the minimum value is 1ms.
		ReadWriteTimeout: 5 * time.Second,

		// InfoLog is used to log info messages to, but can be nil.
		InfoLog: log.New(os.Stdout, "INFO: ", 0),
		// ErrorLog is used to log error messages to, but can be nil.
		ErrorLog: log.New(os.Stderr, "ERROR: ", 0),
	}

	pool.putTokens = make(chan *bt.Put, len(urls))

	for _, url := range urls {
		producer, err := bt.NewProducer(url, pool.putC, options)
		if err != nil {
			return nil, err
		}

		pool.producers = append(pool.producers, producer)
		pool.putTokens <- bt.NewPut(pool.putC, options)
	}

	for _, producer := range pool.producers {
		producer.Start()
	}

	return pool, nil
}

// Publish 发送延迟消息
func (pool *ProducerPool) Publish(
	topic string,
	message interface{},
	delay int64,
) (uint64, error) {

	messageStr, _ := JsonEncode(message)
	R("发送延迟消息: "+topic+", 消息内容: "+Byte2String(messageStr)+", 延迟:"+Int642String(delay)+" 秒", "")

	putParams := &bt.PutParams{
		Priority: 1024,
		Delay:    time.Duration(delay) * time.Second,
		TTR:      30 * time.Second,
	}

	// params = putParams
	// pool := <-p.Pool

	R(pool, "p2")

	routing := strings.Split(topic, "/")
	tube := routing[1]

	msg := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}

	messageData, err := JsonEncode(msg)
	if err != nil {
		graylog.Info(err)
		return 0, err
	}
	// R(tube, "tube")
	// R(Byte2String(messageData), "messageData")

	put := <-pool.putTokens
	id, err := put.Request(tube, messageData, putParams)
	pool.putTokens <- put

	return id, err
}

// Stop shuts down all the producers in the pool.
func (pool *ProducerPool) Stop() {
	pool.stopOnce.Do(func() {
		for i, producer := range pool.producers {
			producer.Stop()
			pool.producers[i] = nil
		}

		pool.producers = []*bt.Producer{}
	})
}

/* func (p Producer) Conn() Producer {
	pool, err := beanstalk.NewProducerPool([]string{p.link}, p.options)
	if err != nil {
		graylog.Err("Unable to create beanstalk producer pool: " + err.Error())
	}
	defer pool.Stop()

	p.Pool = make(chan *beanstalk.ProducerPool)
	p.Pool <- pool
	return p
} */
