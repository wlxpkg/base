/*
 * @Author: qiuling
 * @Date: 2019-04-30 11:50:59
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-04-30 16:31:50
 */
package pkg

import (
	"github.com/gemnasium/logrus-graylog-hook"
	"github.com/gookit/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	var graylog_ip = config.GetEnv("GRAYLOG_IP", "192.168.3.3")
	var graylog_port = config.GetEnv("GRAYLOG_PORT", "3012")

	hook := graylog.NewAsyncGraylogHook(graylog_ip+":"+graylog_port, map[string]interface{}{})
	defer hook.Flush()
	log.AddHook(hook)
	// R(graylog_ip+":"+graylog_port, "inilog")
}

func Info(data interface{}) {
	log.Info(data)
}

func Err(err interface{}) {
	log.Error(err)
}
