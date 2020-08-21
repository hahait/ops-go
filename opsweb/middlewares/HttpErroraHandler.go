package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops.was.ink/opsweb/utils"
)

// 全局 http 错误处理
func HttpErrorHandler(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(utils.Errors); ok {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": e.Code,
					"msg": e.Msg,
					"errmsg": e.Errmsg,
				})
			} else {
				//panic(r)
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 1,
					"msg": "From HttpErrorHandler 出错啦",
					"errmsg": r,
				})
			}
			c.Abort()
		}
	}()
	c.Next()
}