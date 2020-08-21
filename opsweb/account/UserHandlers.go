package account

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"ops.was.ink/opsweb/utils"
)

// 创建用户
func createUserHandler(c *gin.Context) {
	var u User
	// 数据绑定
	utils.ErrorHandler(c.ShouldBindJSON(&u) ,1, "请求数据绑定失败")
	// 数据验证
	utils.ErrorHandler(userFieldsValidator(&u, false) ,1, "模型字段验证失败")
	// 数据写入数据库
	u.Password = utils.EncryptionPassword(u.Password)
	utils.ErrorHandler(utils.Db.Create(&u).Error ,1, "数据写入数据库失败")

	c.JSON(http.StatusCreated, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("用户: %s 创建成功", u.Name),
	})
}

// 删除用户
type DeleteUser struct {
	ID uint `json:"id"`
}
func deleteUserHandler(c *gin.Context) {
	var (
		du DeleteUser
		u User
	)

	// 数据绑定
	utils.ErrorHandler(c.ShouldBindJSON(&du) ,1, "从前端获取用户ID失败")

	// 确认数据库中存在这个用户
	if errors.Is(utils.Db.Select("name").Find(&u, du.ID).Error, gorm.ErrRecordNotFound) {
		panic(fmt.Sprintf("数据库中未查到该用户, 用户ID: %d", du.ID))
	}
	// 删除用户
	utils.ErrorHandler(utils.Db.Delete(&User{}, du.ID).Error ,1, "数据库中删除该用户失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("用户: %s 删除成功", u.Name),
	})
}

// 更新用户信息
func updateUserHandler(c *gin.Context) {
	var (
		up User
		u User
	)
	utils.ErrorHandler(c.ShouldBindJSON(&up) ,1, "从前端获取用户信息失败")
	utils.ErrorHandler(utils.Db.Find(&u, up.ID).Error ,1, fmt.Sprintf("数据库中查找该用户 %s 失败",u.Name))
	// 这里忽略 password 字段的验证
	utils.ErrorHandler(userFieldsValidator(&up, false,"Password", "IsAdmin", "IsActive") ,1, "模型字段验证失败")
	// 这里忽略 password 字段的更新, 后续单独的函数处理
	utils.ErrorHandler(utils.Db.Omit("password", "is_admin", "is_active").Save(&up).Error ,1, "数据库保存失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("用户: %s 信息修改成功", up.Name),
	})
}

// 查询用户列表, 包含了关联关系的查询
// 定义 filters
type QueryUserFilter struct {
	ID uint `form:"id"`
	Name string `form:"name"`
	Phone string `form:"phone"`
	utils.Pagination
}
func queryUserHandler(c *gin.Context) {
	var (
		u []*User
		qu QueryUserFilter
		total int64
	)

	db := utils.Db.Preload("Group")
	utils.ErrorHandler(c.ShouldBindQuery(&qu), 1, "查询字符串绑定失败")
	if qu.ID != 0 {
		db = db.Where("id = ?", qu.ID)
	}
	if qu.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", qu.Name))
	}
	if qu.Phone != "" {
		db = db.Where("phone LIKE ?", fmt.Sprintf("%%%s%%", qu.Phone))
	}
	if pg, ok := qu.Pagination.CheckPage(); ok {
		db = pg.Paginate(db)
	}

	utils.ErrorHandler(db.Find(&u).Count(&total).Error, 1, "从数据库中查找组列表失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"total": total,
		"results": u,
	})
}

// 查询用户信息; 因为 vue 前段在加载页面时, 访问此路由，以获取用户信息；此时不能以 id 的形式去获取，只能通过 context 的方式获取
func queryUserInfoHandler(c *gin.Context) {
	var (
		u User
	)
	name := c.MustGet("username")
	utils.ErrorHandler(utils.Db.Preload("Group").Preload("Group.Permissions").Where("name = ?", name).Find(&u).Error, 1, "从数据库中查询用户信息失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"results": u,
	})
}

// 更新用户密码
func updateUserPasswordHandler(c *gin.Context) {
	var (
		u User
		upu UserPasswordUpdate
	)
	utils.ErrorHandler(c.ShouldBindJSON(&upu) ,1, "从前端获取值失败")
	utils.ErrorHandler(utils.Db.Find(&u, upu.ID).Error ,1, "从数据库中获取该用户失败")
	utils.ErrorHandler(userFieldsValidator(&upu, false) ,1, "验证数据失败")
	utils.ErrorHandler(utils.Db.Model(&u).Updates(User{Password: utils.EncryptionPassword(upu.Password)}).Error ,1, "数据库更新失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("用户 %s 密码更新成功", u.Name),
	})
}

// 更新用户状态
func updateUserStatusHandler(c *gin.Context) {
	var (
		u User
		uus User
	)
	utils.ErrorHandler(c.ShouldBindJSON(&uus) ,1, "从前端获取值失败")
	utils.ErrorHandler(utils.Db.Find(&u, uus.ID).Error ,1, "从数据库中获取该用户失败")
	utils.ErrorHandler(userFieldsValidator(&uus, true, "IsActive") ,1, "验证数据失败")
	utils.ErrorHandler(utils.Db.Model(&u).Updates(User{IsActive: uus.IsActive}).Error ,1, "数据库更新失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("用户 %s 状态更新成功", u.Name),
	})
}

// 更新用户是否是管理员
func updateUserAdminHandler(c *gin.Context) {
	var (
		u User
		uau User
	)
	utils.ErrorHandler(c.ShouldBindJSON(&uau) ,1, "从前端获取值失败")
	utils.ErrorHandler(utils.Db.Find(&u, uau.ID).Error ,1, "从数据库中获取该用户失败")
	utils.ErrorHandler(userFieldsValidator(&uau, true, "IsAdmin") ,1, "验证数据失败")
	utils.ErrorHandler(utils.Db.Model(&u).Updates(User{IsAdmin: uau.IsAdmin}).Error ,1, "数据库更新失败")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("用户 %s 管理员更新成功", u.Name),
	})
}

// 更新用户所属组
type UpdateUserGroups struct {
	ID uint `json:"id"`
	Groups []uint `json:"groups"`
}
func updateUserGroupsHandler(c *gin.Context) {
	var (
		u User
		uug UpdateUserGroups
		gs []Group
	)
	utils.ErrorHandler(c.ShouldBindJSON(&uug), 1, "从前端获取关联的User及Group 对象失败")
	utils.ErrorHandler(utils.Db.Find(&u, uug.ID).Error, 1, "从数据库中获取 User ID 失败")
	if len(uug.Groups) == 0 {
		utils.ErrorHandler(fmt.Errorf("出错啦"), 1, "用户组不能为空")
	}
	utils.ErrorHandler(utils.Db.Where("id IN (?)", uug.Groups).Find(&gs).Error, 1, "数据库查询用户组失败")
	utils.ErrorHandler(utils.Db.Model(&u).Association("Group").Replace(&gs), 1, fmt.Sprintf("数据库更新用户 %s 所属组失败", u.Name))
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg": fmt.Sprintf("更新用户关联组成功", u.Name),
	})
}
