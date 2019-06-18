package pkg

import "errors"

type Errors struct {
	Code    int
	Message string
}

var Errs = map[string]Errors{
	"ERR_UNKNOW_ERROR": Errors{0, "未知系统错误"},
	"ERR_NOERROR":      Errors{1, ""},
	"SUCCESS":          Errors{1, ""},

	"ERR_PARAM":         Errors{10100, "参数错误"},
	"ERR_INVALID_TOKEN": Errors{10102, "无效的token"},
	"ERR_UNAUTHORIZED":  Errors{10104, "您没有权限访问该数据"},
	"ERR_TCP_TIMEOUT":   Errors{10504, "TCP接口响应超时"},
	"ERR_HTTP_TIMEOUT":  Errors{10505, "HTTP接口响应超时"},

	"ERR_MYSQL":              Errors{20100, "MySQL错误"},
	"ERR_REDIS":              Errors{20101, "Reids错误"},
	"ERR_MYSQL_INSTALL_FAIL": Errors{20102, "数据插入失败"},
	"ERR_MYSQL_DELETE_FAIL":  Errors{20103, "数据删除失败"},
	"ERR_IDGEN_FAIL":         Errors{20404, "id生成失败"},
	"ERR_MYSQLPOOL_FAIL":     Errors{21404, "mysql连接池丢失"},
	"ERR_REDISPOOL_FAIL":     Errors{22404, "redis连接池丢失"},

	"ERR_VIDEO_NOT_EXIST":        Errors{30100, "视频不存在"},
	"ERR_VIDEO_EXCEEDED_MAXIMUM": Errors{30200, "翻页总条数超过最大限制"},
	"ERR_GET_FAIL":               Errors{32304, "获取视频异常"},
}

func Excp(errString string) error {
	return errors.New(errString)
}
