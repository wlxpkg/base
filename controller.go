package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	Ctx *gin.Context
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

func (c *Controller) Error(e Errors) {
	c.Ctx.JSON(http.StatusOK, gin.H{
		"code":    e.Code,
		"message": e.Message,
		"data":    "",
	})
	c.Ctx.Abort()
	// return nil
}
