package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ops.was.ink/opsweb/utils"
)

func PermissionHadler(c *gin.Context) {
	rpath := c.Request.URL.String()
	rmethod := c.Request.Method
	if rpath == "/login" {
		c.Next()
	} else {
		ruser := c.MustGet("username")
		if err := utils.Efc.LoadPolicy(); err != nil {
			fmt.Println("load policy 出错, 错误信息: ", err.Error())
		}

		//results, err := utils.Efc.GetRolesForUser("14")
		//if err != nil {
		//	fmt.Println("执行 GetRolesForUser 遇到了错误, 错误信息: ", err)
		//}
		//fmt.Println("PermissionHandler 获取到的 roles: ", results)
		//
		//rm_obj := utils.Efc.GetRoleManager()
		//aa, _ := rm_obj.GetRoles("14")
		//fmt.Println("我此时查询得到的 roles : ", aa)

		b, err := utils.Efc.Enforce(ruser, rpath, rmethod)
		if err != nil {
			utils.ErrorHandler(err, 1, "权限检查出现错误")
		}
		if !b {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 1,
				"msg": "你无权访问",
				"errmsg": "权限验证失败",
			})
			c.Abort()
		}
		c.Next()
	}

	//rm_obj := utils.Efc.GetRoleManager()
	//aa, _ := rm_obj.GetRoles("14")
	//fmt.Printf("获取到的 rm 对象是: %#v \n", aa)
}
