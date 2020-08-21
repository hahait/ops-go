package account

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ops.was.ink/opsweb/utils"
)

// 查询权限
// 包含过滤和分页
type PermissionFilter struct {
	ID uint `form:"id"`
	utils.Pagination
	Path string `form:"path"`
	Method string `form:"method"`
}
func queryPermissionHandler(c *gin.Context) {
	var (
		p []*Permissions
		total int64
		pf PermissionFilter
	)
	db := utils.Db.Model(&Permissions{}).Order("path desc")
	utils.ErrorHandler(c.ShouldBindQuery(&pf), 1, "查询字符串绑定失败")
	if pf.ID != 0 {
		db = db.Where("id = ?", pf.ID)
	}
	if pf.Path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", pf.Path))
	}
	if pf.Method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", pf.Method))
	}
	db = db.Count(&total)
	if pg, ok := pf.Pagination.CheckPage(); ok {
		db = pg.Paginate(db)
	}
	utils.ErrorHandler(db.Find(&p).Error, 1, "查询权限列表出错")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"total": total,
		"results": p,
	})
}

// 修改权限描述信息
func updatePermissionDescribeHandler(c *gin.Context) {
	var (
		upd Permissions
		p Permissions
	)
	utils.ErrorHandler(c.ShouldBindJSON(&upd), 1, "从前端接收值失败")
	fmt.Println("我此时获取到的权限信息: ", upd)
	utils.ErrorHandler(utils.Db.Find(&p, upd.ID).Error, 1, "从数据库中查询权限失败")
	utils.ErrorHandler(utils.Db.Where("id = ?", upd.ID).Updates(&Permissions{Describe: upd.Describe}).Error, 1, "更新权限描述失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("权限: %s: %s 描述更新成功", p.Method, p.Path),
	})
}