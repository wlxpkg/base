package pkg

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserID   int64  `json:"user_id"`
	Code     string `json:"code"`
	Phone    string `json:"phone"`
	Jwt      string `json:"jwt"`
	Avatar   string `json:"avatar"`
	Nickname string `json:"nickname"`
	Pid      int64  `json:"pid"`
}

type Middleware struct {
	Permission int64
	Token      string
	UserID     int64
	UserInfo   UserInfo
}

type Controller struct {
	Ctx        *gin.Context
	UserID     int64
	UserInfo   UserInfo
	Token      string
	Permission int64
	Jwt        string
	Client     string
	ClientID   string
	AppName    string
	AppVersion string
}

func NewController(ctx *gin.Context) (ctl *Controller) {
	ctl = &Controller{Ctx: ctx}
	ctl.getLoginInfo()
	ctl.getHeaders()
	return
}

// 必须通过中间件拿数据
// 要么后台的用 casbin
// 要么前台的用 member
func (ctl *Controller) getLoginInfo() {
	c := ctl.Ctx

	path := c.Request.URL.Path
	routeSli := strings.Split(path, "/")
	// R(routeSli, "routeSli")

	// 初始值
	ctl.Permission = 0
	ctl.UserInfo = UserInfo{}
	ctl.Token = ""
	ctl.UserID = 0

	if routeSli[1] == "login" || routeSli[1] == "callback" {
		return
	}

	if m, ok := c.Get("middleware"); ok && m != nil {
		middleware, _ := m.(Middleware)
		// R(middleware, "middleware")
		ctl.UserID = middleware.UserID
		ctl.UserInfo = middleware.UserInfo
		ctl.Token = middleware.Token
		ctl.Permission = middleware.Permission
	}
	return
}

func (ctl *Controller) getHeaders() {
	c := ctl.Ctx

	authorization := c.GetHeader("authorization")
	jwt := strings.TrimPrefix(authorization, "Bearer ")

	ctl.Jwt = jwt
	ctl.Client = c.GetHeader("client")
	ctl.ClientID = c.GetHeader("client-id")
	ctl.AppName = c.GetHeader("app-name")
	ctl.AppVersion = c.GetHeader("version")
}

func (c *Controller) Get(param string) string {
	value := c.Ctx.Query(param)
	return value
}

func (c *Controller) Getd(param string, defaultValue string) string {
	value := c.Ctx.DefaultQuery(param, defaultValue)
	return value
}

func (c *Controller) Post(param string) string {
	value := c.Ctx.Request.PostFormValue(param)
	return value
}

func (c *Controller) Postd(param string, defaultValue string) string {
	value := c.Ctx.DefaultPostForm(param, defaultValue)
	return value
}

// func (c *Controller) Json(json interface{}) (json interface{}, error) {
// 	// var e = new(Errors)
// 	if err := c.Ctx.ShouldBindJSON(&json); err != nil {
// 		log.Warn(err)
// 		c.Error(Errs["ERR_PARAM"])
// 		return
// 	}

// 	return
// }

func (c *Controller) Success(result interface{}) {
	c.Ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "",
		"data":    result,
	})
}

func (c *Controller) Error(e error) {
	errors, ok := Errs[e.Error()]
	if !ok {
		errors = Errs["ERR_UNKNOW_ERROR"]
	}

	c.Ctx.JSON(http.StatusOK, gin.H{
		"code":    errors.Code,
		"message": errors.Message,
		"data":    "",
	})
	c.Ctx.Abort()
}
