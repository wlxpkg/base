/*
 * @Author: qiuling
 * @Date: 2019-04-30 11:50:59
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-09 16:15:34
 */
package log

import (
	graylog "github.com/gemnasium/logrus-graylog-hook"
	"github.com/gookit/config"

	// "github.com/pkg/errors"
	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	var graylog_ip = config.GetEnv("GRAYLOG_IP", "192.168.3.3")
	var graylog_port = config.GetEnv("GRAYLOG_PORT", "3012")

	hookStack := logrus_stack.StandardHook()
	hook := graylog.NewGraylogHook(graylog_ip+":"+graylog_port, map[string]interface{}{})
	// NewAsyncGraylogHook NewGraylogHook
	// defer hook.Flush()
	log.AddHook(hookStack)
	log.AddHook(hook)
	// R(graylog_ip+":"+graylog_port, "inilog")
}

func Info(data interface{}) {
	log.Info(data)
}

func Err(err interface{}) {
	// errs := error(errors.New(fmt.Sprint(err)))
	// errs = errors.Wrap(errs, "log Err")
	log.Error(err)
}

func Debug(data interface{}) {
	log.Debug(data)
}

func Warn(err interface{}) {
	// errs := error(errors.New(fmt.Sprint(err)))
	// errs = errors.Wrap(errs, "log Err")
	log.Warn(err)
}

func Panic(data interface{}) {
	log.Panic(data)
}
