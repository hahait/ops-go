package main

//** 定义命名标准 **//
// 结构体名称: 所有单词首字母都大写
// 内部函数名: 驼峰结构（第一个单词首字母小写，后续单词首字母大写）
// 外部函数名和变量名: 同结构体名称
// 内部变量名：全部小写, 单词间以下划线连接

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ops.was.ink/opsweb/account"
	"ops.was.ink/opsweb/grouters"
	mw "ops.was.ink/opsweb/middlewares"
	"ops.was.ink/opsweb/utils"
	"strings"
)


func main() {
	e := gin.Default()
	Jwt := mw.JwtAuthTokenInit()

	e.Use(
		mw.HttpErrorHandler,
		Jwt.MiddlewareFunc(),
		mw.PermissionHadler,
	)
	e.POST("/login", Jwt.LoginHandler)
	e.POST("/logout", Jwt.LogoutHandler)
	e.PUT("/token/refresh", Jwt.RefreshHandler)

	grouters.RouterInit(e)

	// 自动注册路由到权限表中
	var (
		p account.Permissions
		s []string
	)
	routes := e.Routes()
	for _, i := range routes {
		if i.Path != "/login" && i.Path != "/logout" && i.Path != "/token/refresh" {
			s = strings.Split(i.Handler,".")
			d := strings.TrimSuffix(s[len(s)-1], "Handler")
			if num := utils.Db.Raw("select id from permissions where path = ? and method = ?", i.Path, i.Method).Find(&p).RowsAffected; num == 0{
				utils.Db.Create(&account.Permissions{Path: i.Path, Method: i.Method, Describe: d})
			}
		}
	}

	if err := e.Run(":10080"); err != nil {
		fmt.Println("startup service failed, err:%v\n", err)
	}
}
