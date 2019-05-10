package pkg

type Errors struct {
	Code    int
	Message string
}

var Errs = map[string]Errors{
	"ERR_UNKNOW_ERROR": Errors{0, "未知系统错误"},
	"ERR_NOERROR":      Errors{1, ""},
	"SUCCESS":          Errors{1, ""},

	"ERR_PARAM":         Errors{10100, "参数错误"},
	"ERR_INVALID_TOKEN": Errors{10102, "非法的token"},
	"ERR_TCP_TIMEOUT":   Errors{10504, "TCP接口响应超时"},
	"ERR_HTTP_TIMEOUT":  Errors{10505, "HTTP接口响应超时"},

	"ERR_MYSQL":              Errors{20100, "MySQL错误"},
	"ERR_MYSQL_INSTALL_FAIL": Errors{20102, "数据插入失败"},
	"ERR_MYSQL_DELETE_FAIL":  Errors{20103, "数据删除失败"},
	"ERR_IDGEN_FAIL":         Errors{20404, "id生成失败"},
	"ERR_MYSQLPOOL_FAIL":     Errors{21404, "mysql连接池丢失"},
	"ERR_REDISPOOL_FAIL":     Errors{22404, "redis连接池丢失"},
}
