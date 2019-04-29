package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	ctx gin.Context
}

func (c *Controller) Get(param string) string {
	value := c.ctx.Query(param)
	return value
}

func (c *Controller) Getd(param string, defaultValue string) string {
	value := c.ctx.DefaultQuery(param, defaultValue)
	return value
}

func (c *Controller) Post(param string) string {
	value := c.ctx.Request.PostFormValue(param)
	return value
}

func (c *Controller) Postd(param string, defaultValue string) string {
	value := c.ctx.DefaultPostForm(param, defaultValue)
	return value
}

func (c *Controller) Json(json interface{}) interface{} {
	// var e = new(Errors)
	if err := c.ctx.ShouldBindJSON(&json); err != nil {
		return c.Error(ERR_NOERROR)
	}

	return json
}

func (c *Controller) Success(result interface{}) {
	c.ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "",
		"data":    result,
	})
}

func (c *Controller) Error(e Errors) string {
	c.ctx.JSON(http.StatusOK, gin.H{
		"code":    e.Code,
		"message": e.Message,
		"data":    "",
	})
	return ""
}
