package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"ops.was.ink/opsweb/utils"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"strconv"
)

// 创建组
func createGrouphandler(c *gin.Context) {
	var g Group
	utils.ErrorHandler(c.ShouldBindJSON(&g) ,1, "从前端获取数据失败")
	utils.ErrorHandler(groupFieldsValidator(&g, false) ,1, "数据验证失败")
	utils.ErrorHandler(utils.Db.Create(&g).Error ,1, "数据库保存失败")
	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("组 %s 创建成功", g.Name),
	})
}

// 删除组
type DelGroup struct {
	ID uint
}
func deleteGroupHandler(c *gin.Context) {
	var (
		dg DelGroup
		g Group
	)
	utils.ErrorHandler(c.ShouldBindJSON(&dg), 1, "从前端获取组ID失败")
	utils.ErrorHandler(utils.Db.Find(&g, dg.ID).Error, 1, "从前端获取组ID失败")
	utils.ErrorHandler(utils.Db.Where("id = ?", dg.ID).Delete(Group{}).Error, 1, "数据库执行失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("组 %s 删除成功", g.Name),
	})
}

// 更新组
func updateGroupHandler(c *gin.Context) {
	var (
		g Group
		gu Group
	)
	utils.ErrorHandler(c.ShouldBindJSON(&gu), 1, "从前端获取数据失败")
	utils.ErrorHandler(groupFieldsValidator(&gu, false), 1, "数据验证失败")
	utils.ErrorHandler(utils.Db.Find(&g, gu.ID).Error, 1, fmt.Sprintf("从数据库中获取组ID失败"))
	utils.ErrorHandler(utils.Db.Save(&gu).Error, 1, "数据保存失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("组 %s 更新成功", gu.Name),
	})
}

// 查询组列表;
// 用户组 filters; 支持模糊查询组名和分页
type QueryGroupFilter struct {
	ID uint `form:"id"`
	Name string `form:"name"`
	utils.Pagination
}
func queryGroupHandler(c *gin.Context) {
	var (
		gs []*Group
		qg QueryGroupFilter
		total int64
	)
	db := utils.Db.Preload(clause.Associations)
	utils.ErrorHandler(c.ShouldBindQuery(&qg), 1, "查询字符串绑定失败")
	if qg.ID != 0 {
		db = db.Where("id = ?", qg.ID)
	}
	if qg.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", qg.Name))
	}
	if pg, ok := qg.Pagination.CheckPage(); ok {
		db = pg.Paginate(db)
	}
	utils.ErrorHandler(db.Find(&gs).Count(&total).Error, 1, "从数据库中查找组列表失败")

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"total": total,
		"results": gs,
	})

}

// 更新组内用户
type UpdateGroupUsers struct {
	ID uint `json:"id"`
	Users []uint `json:"users"`
}
func updateGroupUsersHandler(c *gin.Context) {
	var (
		g Group
		ugu UpdateGroupUsers
		us []User
	)
	utils.ErrorHandler(c.ShouldBindJSON(&ugu), 1, "从前端获取关联的Group及User对象失败")
	utils.ErrorHandler(utils.Db.Find(&g, ugu.ID).Error, 1, fmt.Sprintf("从数据库中获取 Group: %d 失败", ugu.ID))
	if len(ugu.Users) != 0 {
		utils.ErrorHandler(utils.Db.Where("id IN (?)", ugu.Users).Find(&us).Error, 1, "数据库查询用户失败")
	}
	utils.ErrorHandler(utils.Db.Model(&g).Association("User").Replace(&us), 1, fmt.Sprintf("数据库更新用户组 %s 关联的失败", g.Name))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("更新用户组 %s 关联用户成功", g.Name),
	})
}

// 查询组内用户
func queryGroupUsersHandler(c *gin.Context){
	var (
		g Group
		us []*User
	)
	gid, exist := c.GetQuery("id")
	if !exist {
		utils.ErrorHandler(errors.New("出错啦"), 1, "未从前端获得组ID")
	}
	if nexist := utils.Db.Find(&g, gid).Error; errors.Is(nexist, gorm.ErrRecordNotFound) {
		utils.ErrorHandler(errors.New("出错啦"), 1, "未从数据库中查到该组ID")
	}
	utils.ErrorHandler(utils.Db.Model(&g).Association("User").Find(&us), 1, fmt.Sprintf("从数据库中查找组 %s 下的用户失败", g.Name))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": us,
	})
}

// 修改组的权限
type UpdateGroupPerms struct {
	ID uint64 `json:"id"`
	Permissions []int64 `json:"permissions"`
}
func updateGroupPermsHandlers(c *gin.Context) {
	var (
		g Group
		ugp UpdateGroupPerms
		plist []Permissions
		croles []*gormadapter.CasbinRule
	)
	utils.ErrorHandler(c.ShouldBindJSON(&ugp), 1, "从前端接收值失败")
	utils.ErrorHandler(utils.Db.Find(&g, ugp.ID).Error, 1, "从数据库中查询组失败")
	if len(ugp.Permissions) != 0 {
		utils.ErrorHandler(utils.Db.Where("id IN (?)", ugp.Permissions).Find(&plist).Error, 1, "从数据库中获取权限失败")
	}
	utils.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&g).Association("Permissions").Replace(&plist); err != nil {
			return errors.New(fmt.Sprintf("更新用户组: %s 权限失败", g.Name))
		}
		if err := tx.Where(gormadapter.CasbinRule{PType: "p", V0: strconv.FormatUint(ugp.ID, 10)}).Delete(gormadapter.CasbinRule{}).Error ; err != nil {
			return errors.New(fmt.Sprintf("从 casbin_rule 中删除组: %s 权限失败", g.Name))
		}

		if len(ugp.Permissions) != 0 {
			for _, i := range plist {
				croles = append(croles, &gormadapter.CasbinRule{
					PType: "p",
					V0: strconv.FormatUint(ugp.ID, 10),
					V1: i.Path,
					V2: i.Method,
				})
			}
			if err := tx.Create(&croles).Error; err != nil {
				return errors.New("向 casbin_rule 表中重新添加权限失败")
			}
		}
		return nil
	})
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("修改组: %s 权限成功", g.Name),
	})
}