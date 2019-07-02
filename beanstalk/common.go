/*
 * @Author: qiuling
 * @Date: 2019-07-02 16:21:23
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-02 16:24:42
 */
package beanstalk

import (
	. "artifact/pkg/config"
	"log"
	"os"
	"time"

	bt "github.com/prep/beanstalk"
)

func GetOptions() (urls []string, options *bt.Options) {
	link := "beanstalk://" + Config.Beanstalk.Host + ":" + Config.Beanstalk.Port

	urls = append(urls, link, link)

	options = &bt.Options{
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
	return
}
