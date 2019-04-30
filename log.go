/*
 * @Author: qiuling
 * @Date: 2019-04-30 11:50:59
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-04-30 15:07:57
 */
package pkg

import (
	"github.com/gemnasium/logrus-graylog-hook"
	"github.com/gookit/config"
	log "github.com/sirupsen/logrus"
)

func init() {
	var graylog_ip = config.GetEnv("GRAYLOG_IP", "192.168.3.3")
	var graylog_port = config.GetEnv("GRAYLOG_UDP_PORT", "12500")

	hook := graylog.NewAsyncGraylogHook(graylog_ip+":"+graylog_port, map[string]interface{}{})
	defer hook.Flush()
	log.AddHook(hook)
	R(graylog_ip+":"+graylog_port, "inilog")
}

func LogInfo(data interface{}) {
	log.Info(data)
}

func LogErr(err interface{}) {
	log.Error(err)
}
