/*
 * @Author: qiuling
 * @Date: 2019-06-28 15:38:14
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-02 16:28:05
 */
package beanstalk

import (
	. "artifact/pkg"
	"artifact/pkg/log"
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

	urls, options := GetOptions()

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
		log.Err("Unable to create beanstalk producer pool: " + err.Error())
	}
	defer pool.Stop()

	p.Pool = make(chan *beanstalk.ProducerPool)
	p.Pool <- pool
	return p
} */
