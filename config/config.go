/*
 * @Author: qiuling
 * @Date: 2019-05-13 16:01:39
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-05-15 14:49:32
 */

package pkg

import (
	"artifact/pkg/log"
	"os"
	"path/filepath"
	"strings"

	"github.com/koding/multiconfig"
)

type (
	// Server holds supported types by the multiconfig package
	/* Conf struct {
		Port        int `default:"8000"`
		Redis       redis
		Mysql       mysql
		Mongo       mongo
		Idgenerator idgenerator
		Amqp        amqp
		Beanstalk   beanstalk
		Graylog     graylog
		Oss         oss
	} */

	redis struct {
		Host     string `default:"192.168.3.3"`
		Port     string `default:"6379"`
		Select   int    `default:"0"`
		Password string
		Prefix   string
	}

	mysql struct {
		Port     string `default:"3306"`
		Host     string `default:"192.168.3.3"`
		Database string
		Username string
		Password string
	}

	mongo struct {
		Port     string `default:"3306"`
		Host     string `default:"192.168.3.3"`
		Database string
		Username string
		Password string
	}

	idgenerator struct {
		Url string
	}

	amqp struct {
		Address string `default:"192.168.3.3"`
		Port    int    `default:"5672"`
		User    string `default:"artifact"`
		Pass    string `default:"artifact"`
		Vhost   string `default:"/artifact"`
	}

	beanstalk struct {
		Addr string `default:"192.168.3.3"`
		Port int    `default:"11300"`
	}

	graylog struct {
		Host string `default:"192.168.3.3"`
		Port int    `default:"3012"`
	}

	oss struct {
		Id     string
		Secret string
		Host   string
		Bucket string
		Sts    ossSts
	}

	ossSts struct {
		Id     string
		Secret string
		Arn    string
	}
)

var Config = new(struct {
	Port        int `default:"8000"`
	Redis       redis
	Mysql       mysql
	Mongo       mongo
	Idgenerator idgenerator
	Amqp        amqp
	Beanstalk   beanstalk
	Graylog     graylog
	Oss         oss
})

func init() {
	path := getCurrentDirectory()
	parent := getParentDirectory(path)

	m := multiconfig.NewWithPath(parent + "/config/" + "config.toml")
	m.MustLoad(Config)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Info(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
