/*
 * @Author: qiuling
 * @Date: 2019-06-28 15:38:14
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-06-28 23:28:39
 */
package beanstalk

import (
	. "artifact/pkg"
	. "artifact/pkg/config"
	graylog "artifact/pkg/log"
	"log"
	"os"
	"strings"
	"time"

	bt "github.com/prep/beanstalk"
)

type Producer struct {
	bt.ProducerPool
	link    string
	options *bt.Options
	// Pool    chan *bt.ProducerPool
}

func NewProducer() *Producer {
	p := new(Producer)
	// R(Config.Beanstalk.Host, "host")
	// R(Config.Beanstalk.Port, "Port")

	link := "beanstalk://" + Config.Beanstalk.Host + ":" + Config.Beanstalk.Port

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

	// pool, err := beanstalk.NewProducerPool([]string{link}, options)
	// if err != nil {
	// 	graylog.Err("Unable to create beanstalk producer pool: " + err.Error())
	// }
	// defer pool.Stop()
	p.link = link
	p.options = options

	/* pool, err := beanstalk.NewProducerPool([]string{p.link}, p.options)
	if err != nil {
		graylog.Err("Unable to create beanstalk producer pool: " + err.Error())
	}
	defer pool.Stop()

	p.Pool = make(chan *beanstalk.ProducerPool)
	p.Pool <- pool

	R(p, "p0")
	defer close(p.Pool) */
	return p
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

func (p Producer) NewPool() (*bt.ProducerPool, error) {
	return bt.NewProducerPool([]string{p.link, p.link, p.link, p.link, p.link, p.link}, p.options)
}

func (p Producer) MakeData(topic string, message interface{}, delay int64) (tube string, messageData []byte, putParams *bt.PutParams, err error) {
	putParams = &bt.PutParams{
		Priority: 1024,
		Delay:    time.Duration(delay) * time.Second,
		TTR:      30 * time.Second,
	}

	// params = putParams
	// pool := <-p.Pool

	// R(pool, "p2")

	routing := strings.Split(topic, "/")
	tube = routing[1]

	msg := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}

	messageData, err = JsonEncode(msg)
	if err != nil {
		graylog.Info(err)
		return
	}
	R(tube, "tube")
	R(Byte2String(messageData), "messageData")

	// jobID = 111
	// err = Excp("!11")
	// return

	/* // jobID, err = p.pool.Put(tube, messageData, putParams)
	jobID, err = pool.Put("default", []byte("Hello World"), putParams)
	p.Pool <- pool
	R(jobID, "jobID")
	R(err, "err")
	if err != nil {
		return
	}
	return */

	return
}
