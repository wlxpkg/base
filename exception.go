package pkg

import (
	"io"
	"runtime"

	"artifact/pkg/log"

	"github.com/gin-gonic/gin"
)

type Exception struct {
}

// Error Struct
/* type Error struct {
	StatusCode int
	Msg        string
	Code       int
}

func (err *Error) Error() string {
	return fmt.Sprintf("status_code:%d, msg:%s", err.StatusCode, err.Msg)
}

// RegisterErrors : register your error messages
func RegisterErrors(msgFlags map[int]string) {
	for key, value := range msgFlags {
		ginErrors.MsgFlags[key] = value
	}
}
*/
// GenError error build, You can implement this function yourself.
/* func GenError(httpCode int, errCode int, msg ...string) Error {
	err := Error{
		StatusCode: httpCode,
		Code:       errCode,
	}
	if len(msg) > 0 { // your message stri
		err.Msg = msg[0]
	} else {
		err.Msg = ginErrors.GetMsg(errCode)
	}
	return err
} */

/* func abortWithError(c *gin.Context, err Error) {
	c.JSON(err.StatusCode, gin.H{
		"code":    err.Code,
		"message": err.Msg,
	})
	c.Abort()
} */

// Stack returns a formatted stack trace of the goroutine that calls it.
// It calls runtime.Stack with a large enough buffer to capture the entire trace.
func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// ErrorHandle 统一捕获错误
func ErrorHandle(out io.Writer) gin.HandlerFunc {
	// logger := log.New(out, "", log.LstdFlags|log.Llongfile)
	return func(ctx *gin.Context) {
		defer func() {
			ctl := Controller{Ctx: ctx}

			if err := recover(); err != nil {
				if e, ok := err.(Errors); ok {
					//自定义错误，业务逻辑故意抛出的，返回统一格式数据
					ctl.Error(e)
					// abortWithError(ctx, e)
					return
				}

				// error stack
				stack := make([]byte, 1024*8)
				stack = stack[:runtime.Stack(stack, false)]
				log.Err(Stack())
				// httprequest, _ := httputil.DumpRequest(ctx.Request, false)
				// logger.Printf("[Handler] panic recovered:\n%s\n%s\n%s", string(httprequest), err, Stack())
				// default error
				// internalServerError := GenError(http.StatusInternalServerError, ginErrors.ERROR)
				// abortWithError(ctx, internalServerError)
				ctl.Error(ERR_UNKNOW_ERROR)
			}
		}()
		ctx.Next()
	}
}
