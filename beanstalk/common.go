/*
 * @Author: qiuling
 * @Date: 2019-07-02 16:21:23
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-08-08 18:08:18
 */
package beanstalk

import (
	. "artifact/pkg/config"
	"log"
	"os"
	"time"

	bt "github.com/prep/beanstalk"
)

func GetOptions() (urls []string, config bt.Config) {
	link := "beanstalk://" + Config.Beanstalk.Host + ":" + Config.Beanstalk.Port

	urls = append(urls, link, link)

	config = bt.Config{
		// ReserveTimeout is the time a consumer should wait before reserving a job,
		// when the last attempt didn't yield a job.
		// The default is to wait 5 seconds.
		ReserveTimeout: 5 * time.Second,
		// ReconnectTimeout is the timeout between reconnects.
		// The default is to wait 10 seconds.
		ReconnectTimeout: 10 * time.Second,

		// InfoLog is used to log info messages to, but can be nil.
		InfoLog: log.New(os.Stdout, "INFO: ", 0),
		// ErrorLog is used to log error messages to, but can be nil.
		ErrorLog: log.New(os.Stderr, "ERROR: ", 0),
	}
	return
}
