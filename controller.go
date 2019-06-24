package pkg

import (
	// . "artifact/pkg/middleware"
	"artifact/pkg/biz"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserInfo struct {
	UserID   int64
	Code     string
	Phone    string
	Jwt      string
	Avator   string
	Nickname string
	Pid      int64
}

type Middleware struct {
	Permission int64
	Token      string
	UserID     int64
	// Hostname   string
	UserInfo UserInfo
}

type Controller struct {
	Ctx        *gin.Context
	UserID     int64
	UserInfo   UserInfo
	Token      string
	Permission int64
	Jwt        string
	ClientID   string
	AppName    string
	AppVersion string
}

func NewController(ctx *gin.Context) (ctl Controller) {
	ctl = Controller{Ctx: ctx}
	ctl = getLoginInfo(ctx, ctl)
	ctl = getHeaders(ctx, ctl)
	return
}

func getLoginInfo(c *gin.Context, ctl Controller) Controller {

	if m, ok := c.Get("middleware"); ok && m != nil {
		middleware, _ := m.(Middleware)
		// R(middleware, "middleware")
		ctl.UserID = middleware.UserID
		ctl.UserInfo = middleware.UserInfo
		ctl.Token = middleware.Token
		ctl.Permission = middleware.Permission
	} else {
		token, userInfo, err := GetUser(c)

		if err != nil {
			ctl.Error(err)
		}

		pid, _ := String2Int64(userInfo["pid"])
		uid, _ := String2Int64(userInfo["user_id"])

		info := UserInfo{
			UserID:   uid,
			Code:     userInfo["code"],
			Phone:    userInfo["phone"],
			Jwt:      userInfo["jwt"],
			Avator:   userInfo["avator"],
			Nickname: userInfo["nickname"],
			Pid:      pid,
		}

		ctl.Permission = 0
		ctl.UserInfo = info
		ctl.Token = token
		ctl.UserID = uid
	}

	return ctl
}

func getHeaders(c *gin.Context, ctl Controller) Controller {

	authorization := c.GetHeader("authorization")
	jwt := strings.TrimPrefix(authorization, "Bearer ")

	ctl.Jwt = jwt
	ctl.ClientID = c.GetHeader("client-id")
	ctl.AppName = c.GetHeader("app-name")
	ctl.AppVersion = c.GetHeader("version")

	return ctl
}

func GetUser(c *gin.Context) (token string, userInfo map[string]string, err error) {
	authorization := c.GetHeader("authorization")
	jwt := strings.TrimPrefix(authorization, "Bearer ")

	// fmt.Printf("jwt:%+v\n", jwt)

	if jwt == "" {
		err = errors.New("ERR_INVALID_TOKEN")
		return
	}

	token, err = biz.Jwt2Token(jwt)

	if token == "" || err != nil {
		err = errors.New("ERR_INVALID_TOKEN")
		return
	}

	userInfo = biz.TokenGetUser(token)
	// fmt.Printf("userInfo:%+v\n", userInfo)
	return
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
