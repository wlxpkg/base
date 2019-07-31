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

	"ERR_PARAM":         Errors{1000, "参数错误"},
	"ERR_UNLOGIN":       Errors{1001, "请先注册登录"},
	"ERR_INVALID_TOKEN": Errors{1002, "无效的token"},
	"ERR_UNAUTHORIZED":  Errors{1004, "您没有权限访问该数据"},
	"ERR_DATA_DECODE":   Errors{1005, "数据解析失败"},
	"ERR_HTTP_BASEURL":  Errors{1006, "请设置 baseurl"},
	"ERR_TCP_TIMEOUT":   Errors{1504, "TCP接口响应超时"},
	"ERR_HTTP_TIMEOUT":  Errors{1505, "HTTP接口响应超时"},

	"ERR_MYSQL":              Errors{2000, "MySQL错误"},
	"ERR_MYSQL_INSTALL_FAIL": Errors{2001, "数据插入失败"},
	"ERR_MYSQL_DELETE_FAIL":  Errors{2002, "数据删除失败"},
	"ERR_MYSQLPOOL_FAIL":     Errors{2004, "mysql连接池丢失"},
	"ERR_REDIS":              Errors{2100, "Reids错误"},
	"ERR_REDISPOOL_FAIL":     Errors{2104, "redis连接池丢失"},
	"ERR_IDGEN_FAIL":         Errors{2404, "id生成失败"},

	"ERR_PAY_FUBEI_FAIL": Errors{20000, "支付请求失败"},

	"ERR_VIDEO_NOT_EXIST":        Errors{25000, "视频不存在"},
	"ERR_VIDEO_EXCEEDED_MAXIMUM": Errors{25001, "翻页总条数超过最大限制"},
	"ERR_GET_FAIL":               Errors{25002, "获取视频异常"},

	"ERR_ORDER_NOT_EXIST":            Errors{21000, "订单不存在"},
	"ERR_ORDER_ALREADY_PAY":          {21001, "订单已经支付"},
	"ERR_ORDER_TIME_OUT":             {21002, "订单已经超时，请重新下单"},
	"ERR_ORDER_RESET_TYPE_NOT_EXIST": {21003, "重置订单的支付类型不存在"},

	"ERR_LINK_APPLY_AGENT_INVALD":        {21101, "申请代理的链接失败"},
	"ERR_LINK_APPLY_AGENT_ALREADY_AGENT": {21102, "申请代理的用户已经是代理用户"},

	"ERR_WALLET_PAY_FAIL": {22000, "钱包支付失败"},
	"ERR_WALLET_PAY_EXP":  {22001, "钱包数额不正确"},
}

func Excp(errString string) error {
	return errors.New(errString)
}
