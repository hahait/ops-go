package utils

import (
	"errors"
	"github.com/casbin/casbin/v2/rbac"
	"gorm.io/gorm"
	"strconv"
)

type RoleManager struct {
	maxHierarchyLevel int
}

func NewRoleManager(maxHierarchyLevel int) rbac.RoleManager {
	rm := RoleManager{
		maxHierarchyLevel: maxHierarchyLevel,
	}
	return &rm
}

// 由于 用户与角色间的映射关系不需要 role manager 管理而是通过 用户与组的关联关系定义；
// 这里的 role manager 仅用来获取用户与角色的映射；
// 所以相关的 AddLink(), DeleteLink() 方法无需定义；

// 要返回 nil, 因为 e.LoadPolicy() 时会调用到
func (rm *RoleManager) Clear() error {
	//return errors.New("Clear not implemented")
	return nil
}

func (rm *RoleManager) AddLink(name1 string, name2 string, domain ...string) error {
	return errors.New("AddLink not implemented")
}

func (rm *RoleManager) DeleteLink(name1 string, name2 string, domain ...string) error {
	return errors.New("DeleteLink not implemented")
}

func (rm *RoleManager) HasLink(uid string, gid string, domian ...string) (bool, error) {
	var (
		user_id, _ =  strconv.Atoi(uid)
		group_id, _ =  strconv.Atoi(gid)
		ugmap = map[string]interface{}{}
	)
	if err := Db.Table("user_groups").Where("user_id = ? and group_id = ?", user_id, group_id).Find(&ugmap).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return false , err
	}
	return true, nil
}

func (rm *RoleManager) GetRoles(uid string, domain ...string) ([]string, error) {
	var (
		user_id, _ =  strconv.Atoi(uid)
		g_map = map[string]interface{}{}
	)
	g_list := make([]string, 0)
	if err := Db.Table("user_groups").Select("group_id").Where("user_id = ?", user_id).Find(&g_map).Error; err != nil {
		return nil, err
	}

	for _, v := range g_map {
		gv, _ := v.(int64)
		g_list = append(g_list, strconv.FormatInt(gv,10))
	}

	return g_list, nil
}

func (rm *RoleManager) GetUsers(gid string, domain ...string) ([]string, error) {
	var (
		group_id, _ =  strconv.Atoi(gid)
		u_map = map[string]interface{}{}
	)
	u_list := make([]string, 0)
	if err := Db.Table("user_groups").Select("user_id").Where("group_id = ?", group_id).Find(&u_map).Error; err != nil {
		return nil, nil
	}
	for _, v := range u_map {
		uv, _ := v.(int64)
		u_list = append(u_list, strconv.FormatInt(uv,10))
	}

	return u_list, nil
}

// 要返回 nil; 因为在 e.LoadPolicy() 时会调用到
func (rm *RoleManager) PrintRoles() error {
	return nil
}
