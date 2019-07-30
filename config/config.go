/*
 * @Author: qiuling
 * @Date: 2019-05-13 16:01:39
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-07-30 15:10:52
 */

package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ql2005/multiconfig"
)

type (
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
		Url  string
		Host string `default:"192.168.3.3"`
		Port string `default:"6389"`
	}

	amqp struct {
		Host  string `default:"192.168.3.3"`
		Port  string `default:"5672"`
		User  string `default:"artifact"`
		Pass  string `default:"artifact"`
		Vhost string `default:"/artifact"`
	}

	beanstalk struct {
		Host string `default:"192.168.3.3"`
		Port string `default:"11300"`
	}

	graylog struct {
		Host string `default:"192.168.3.3"`
		Port string `default:"3012"`
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

	jwt struct {
		Secret string
		Uid    string
	}

	server struct {
		Tools     string
		User      string
		Course    string
		Discovery string
		Common    string
		Grant     string
		Shop      string
		Message   string
		Game      string
	}

	service struct {
		Secret string
	}

	fubei struct {
		Url            string
		AlipayCallback string
		WechatCallback string
	}

	wechat struct {
		Id       string
		Callback string
	}

	alipay struct {
		ReturnUrl string
		NotifyUrl string
	}

	pay struct {
		Fubei  fubei
		Wechat wechat
		Alipay alipay
	}

	rsa struct {
		Public  string
		Private string
	}
)

var Config = new(struct {
	Port        string `default:"8000"`
	Redis       redis
	Mysql       mysql
	Mongo       mongo
	Idgenerator idgenerator
	Amqp        amqp
	Beanstalk   beanstalk
	Graylog     graylog
	Oss         oss
	Jwt         jwt
	Server      server
	Service     service
	Pay         pay
	Rsa         rsa
})

func init() {
	path := GetCurrentDirectory()
	parent := GetParentDirectory(path)

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

func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
		// log.Info(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
