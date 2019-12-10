/*
 * @Author: qiuling
 * @Date: 2019-06-28 15:39:03
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:10:17
 */
package beanstalk

import (
	"context"
	"fmt"
	. "github.com/wlxpkg/base"
	"github.com/wlxpkg/base/log"
	"sync"
	"time"

	bt "github.com/prep/beanstalk"
)

type Callback func(message string) (bool, error)

type Consumer struct {
	wg        sync.WaitGroup
	tube      string
	pool      *bt.ConsumerPool
	receivers []Receiver
}

func NewConsumer(tube string) *Consumer {
	return &Consumer{
		tube: tube,
	}
}

// RegisterReceiver 注册一个用于接收指定队列指定路由的数据接收者
func (c *Consumer) RegisterReceiver(receiver Receiver) {
	c.receivers = append(c.receivers, receiver)
}

func (c *Consumer) Start() {
	for {
		c.run()

		R(time.Now().Format("2006-01-02 15:04:05"), "beanstalk 断开连接")
		// 一旦连接断开，那么需要隔一段时间去重连
		// 这里最好有一个时间间隔
		time.Sleep(3 * time.Second)
	}
}

func (c *Consumer) run() {
	var err error
	urls, config := GetOptions()
	c.pool, err = bt.NewConsumerPool(urls, []string{c.tube}, config)
	R("连接成功", "beanstalk消费者 ")
	if err != nil {
		log.Err("Unable to create beanstalk consumer pool: " + err.Error())
	}
	defer c.pool.Stop()

	for _, receiver := range c.receivers {
		c.wg.Add(1)
		go c.listen(receiver) // 每个接收者单独启动一个goroutine用来初始化queue并接收消息
	}

	c.wg.Wait()

	log.Err("所有处理queue的任务都意外退出了")

	// 理论上c.run()在程序的执行过程中是不会结束的
	// 那么则需要重新连接，这里尝试销毁当前连接
	defer c.distory()

}

func (c *Consumer) listen(receiver Receiver) {
	defer c.wg.Done()

	c.pool.Play()

	ctx := context.Background()

	for job := range c.pool.C {
		// logmsg := fmt.Sprintf("收到延时任务 id: %d body: %s\n", job.ID, string(job.Body))
		// log.Info(logmsg)

		ok := receiver.OnReceive(Byte2String(job.Body))

		if ok {
			_ = job.Delete(ctx)
		} else {
			logmsg := fmt.Sprintf("回退延时任务 id: %d body: %s\n", job.ID, string(job.Body))
			log.Warn(logmsg)
			_ = job.Bury(ctx)
		}
	}
}

func (c *Consumer) distory() {
	go c.pool.Stop()
}
