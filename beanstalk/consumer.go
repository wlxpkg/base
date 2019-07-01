/*
 * @Author: qiuling
 * @Date: 2019-06-28 15:39:03
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-01 10:31:30
 */
package beanstalk

import (
	. "artifact/pkg"
	. "artifact/pkg/config"
	"artifact/pkg/log"
	"fmt"
	syslog "log"
	"os"
	"time"

	"github.com/prep/beanstalk"
	bt "github.com/prep/beanstalk"
)

type Callback func(job *bt.Job) (bool, error)

func NewConsumer(tube string, callback Callback) {
	link := "beanstalk://" + Config.Beanstalk.Host + ":" + Config.Beanstalk.Port

	options := &bt.Options{
		ReserveTimeout:   3 * time.Second,
		ReconnectTimeout: 3 * time.Second,
		ReadWriteTimeout: 5 * time.Second,
		InfoLog:          syslog.New(os.Stdout, "INFO: ", 0),
		ErrorLog:         syslog.New(os.Stderr, "ERROR: ", 0),
	}
	pool, err := beanstalk.NewConsumerPool([]string{link, link}, []string{tube}, options)
	if err != nil {
		log.Err("Unable to create beanstalk consumer pool: " + err.Error())
	}
	defer pool.Stop()

	pool.Play()

	for {
		select {
		case job := <-pool.C:
			R(job, "job")
			log.Info("Received job with id: " + Uint642String(job.ID))

			ok, err := callback(job)

			if ok && err == nil {
				_ = job.Delete()
			} else {
				logmsg := fmt.Sprintf("Burying job %d with body: %s\n", job.ID, string(job.Body))
				log.Warn(logmsg)
			}

		}
	}
}
