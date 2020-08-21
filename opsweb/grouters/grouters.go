package grouters

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ops.was.ink/opsweb/account"
)

type Option func(*gin.Engine)

var options = []Option{
	account.Routers,
}

// 定义一个路由注册初始化函数，注册所有 app 下的路由
func RouterInit(e *gin.Engine) {
	for _, opt := range options {
		opt(e)
	}

	e.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})

}
